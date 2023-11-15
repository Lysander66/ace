package hls

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/Lysander66/ace/pkg/common/cnet"
	"github.com/Lysander66/ace/pkg/playlist"
)

var errInvalidPlaylist = fmt.Errorf("invalid playlist")

type AfterFirstPlaylistDownloadFunc func(playlist.Playlist, []byte)
type ClientAfterDownloadSegmentFunc func(*playlist.MediaSegment, []byte)

type PullSession struct {
	URI                        string
	AfterFirstPlaylistDownload AfterFirstPlaylistDownloadFunc
	AfterDownloadSegment       ClientAfterDownloadSegmentFunc
	numParallel                int
	httpClient                 *cnet.Client
	ctx                        context.Context
	ctxCancel                  func()
	playlistURL                *url.URL
	curMediaSequence           int
	curSegmentID               int
	segmentCh                  chan *playlist.MediaSegment
	outErr                     chan error
}

func NewPullSession(uri string, hc *cnet.Client, afterDownloadSegment ClientAfterDownloadSegmentFunc) *PullSession {
	c := &PullSession{
		URI:                  uri,
		AfterDownloadSegment: afterDownloadSegment,
		numParallel:          1,
		httpClient:           hc,
	}
	return c
}

func (c *PullSession) SetParallel(n int) {
	c.numParallel = n
}

func (c *PullSession) Start() error {
	var err error
	c.playlistURL, err = url.Parse(c.URI)
	if err != nil {
		return err
	}

	if c.httpClient == nil {
		c.httpClient = cnet.New()
	}

	c.ctx, c.ctxCancel = context.WithCancel(context.Background())
	c.segmentCh = make(chan *playlist.MediaSegment, 1)
	c.outErr = make(chan error, 1)

	go c.run()

	return nil
}

// Close closes all the Client resources.
func (c *PullSession) Close() {
	c.ctxCancel()
}

// Wait waits for any error of the Client.
func (c *PullSession) Wait() chan error {
	return c.outErr
}

func (c *PullSession) run() {
	c.outErr <- c.runInner()
}

func (c *PullSession) runInner() error {
	p, body, err := FetchPlaylist(c.httpClient, c.playlistURL.String())
	if err != nil {
		slog.ErrorContext(c.ctx, "FetchPlaylist", "err", err, "URI", c.playlistURL.String())
		return err
	}
	if c.AfterFirstPlaylistDownload != nil {
		c.AfterFirstPlaylistDownload(p, body)
	}

	var initialPlaylist *playlist.Media

	switch pl := p.(type) {
	case *playlist.Media: // Media Playlist
		initialPlaylist = pl
	case *playlist.Multivariant: // Master Playlist
		leadingPlaylist := PickLeadingPlaylist(pl.Variants)
		if leadingPlaylist == nil {
			return fmt.Errorf("no variants with supported codecs found")
		}

		u, err := ClientAbsoluteURL(c.playlistURL, leadingPlaylist.URI)
		if err != nil {
			slog.ErrorContext(c.ctx, "AbsoluteURL", "err", err, "URI", leadingPlaylist.URI)
			return err
		}

		pl2, body, err := FetchMediaPlaylist(c.httpClient, u.String())
		if err != nil {
			slog.ErrorContext(c.ctx, "FetchMediaPlaylist", "err", err, "URI", u.String())
			return err
		}
		if c.AfterFirstPlaylistDownload != nil {
			c.AfterFirstPlaylistDownload(p, body)
		}

		initialPlaylist = pl2
		c.playlistURL = u

		// TODO
		if leadingPlaylist.Audio != "" {
		}
	default:
		return errInvalidPlaylist
	}

	// Loading Media Segments
	for i := 0; i < c.numParallel; i++ {
		go c.loadMediaSegments()
	}

	isVOD := initialPlaylist.Endlist || (initialPlaylist.PlaylistType != nil && *initialPlaylist.PlaylistType == playlist.MediaPlaylistTypeVOD)
	if !isVOD {
		// https://datatracker.ietf.org/doc/html/rfc8216#section-6.3.5
		// Determining the next segment to load
		index := len(initialPlaylist.Segments) - 3
		if index < 0 {
			index = 0
		}

		// Loading Media Segments
		c.curMediaSequence = initialPlaylist.MediaSequence
		for i := index; i < len(initialPlaylist.Segments); i++ {
			c.curSegmentID = initialPlaylist.MediaSequence + i
			c.segmentCh <- initialPlaylist.Segments[i]
		}

		// Reloading the Media Playlist
		c.reloadMediaPlaylist()
	} else {
		for _, seg := range initialPlaylist.Segments {
			c.segmentCh <- seg
		}
	}

	return nil
}

func (c *PullSession) reloadMediaPlaylist() {
	timer := time.NewTimer(0)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			pl, _, err := FetchMediaPlaylist(c.httpClient, c.playlistURL.String())
			if err != nil {
				slog.ErrorContext(c.ctx, "FetchMediaPlaylist", "err", err, "URI", c.playlistURL.String())
				return
			}

			c.curMediaSequence = pl.MediaSequence

			timer.Reset(getInterval(pl.TargetDuration, c.curMediaSequence != pl.MediaSequence))

			for i, seg := range pl.Segments {
				curSegmentID := pl.MediaSequence + i
				if curSegmentID > c.curSegmentID {
					c.curSegmentID = curSegmentID
					c.segmentCh <- seg
				}
			}

			// no more Media Segments
			if pl.Endlist {
				c.ctxCancel()
				slog.Info("reloadMediaPlaylist: end")
			}
		case <-c.ctx.Done():
			slog.Info("reloadMediaPlaylist: quit")
			return
		}
	}
}

func getInterval(targetDuration int, hasChanged bool) time.Duration {
	if hasChanged {
		slog.Debug("getInterval", "interval", targetDuration, "hasChanged", hasChanged)
		return time.Duration(targetDuration) * time.Second
	}

	interval := 1000 * targetDuration / 2
	if interval < 1000 {
		interval = 1000
	}
	slog.Debug("getInterval", "interval", interval, "hasChanged", hasChanged)
	return time.Duration(interval) * time.Millisecond
}

func (c *PullSession) loadMediaSegments() {
	for {
		select {
		case seg := <-c.segmentCh:
			c.downloadSegment(seg)
		case <-c.ctx.Done():
			// 消费完再退出
			for len(c.segmentCh) > 0 {
				seg := <-c.segmentCh
				c.downloadSegment(seg)
			}
			slog.Info("loadMediaSegments: quit")
			return
		}
	}
}

func (c *PullSession) downloadSegment(seg *playlist.MediaSegment) error {
	u, err := ClientAbsoluteURL(c.playlistURL, seg.URI)
	if err != nil {
		slog.Error("downloadSegment", "err", err, "URI", seg.URI)
		return err
	}

	data, err := FetchSegment(c.httpClient, u.String())
	if err != nil {
		slog.Error("FetchSegment", "err", err, "URI", u.String())
		return err
	}

	if c.AfterDownloadSegment != nil {
		c.AfterDownloadSegment(seg, data)
	}

	return nil
}

func FetchMediaPlaylist(httpClient *cnet.Client, rawURL string) (*playlist.Media, []byte, error) {
	pl, body, err := FetchPlaylist(httpClient, rawURL)
	if err != nil {
		return nil, body, err
	}

	plt, ok := pl.(*playlist.Media)
	if !ok {
		return nil, body, errInvalidPlaylist
	}

	return plt, body, nil
}

func FetchPlaylist(httpClient *cnet.Client, rawURL string) (playlist.Playlist, []byte, error) {
	resp, err := httpClient.R().Get(rawURL)
	if err != nil {
		return nil, nil, err
	}
	if resp.IsError() {
		return nil, nil, fmt.Errorf("status: %d", resp.StatusCode())
	}

	pl, err := playlist.Unmarshal(resp.Body())
	if err != nil {
		return nil, nil, err
	}
	return pl, resp.Body(), nil
}

func FetchSegment(httpClient *cnet.Client, rawURL string) ([]byte, error) {
	resp, err := httpClient.R().GetWithRetries(rawURL)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("status: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}
