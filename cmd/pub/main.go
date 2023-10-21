package main

import (
	"context"
	"log/slog"

	"github.com/Lysander66/ace/pkg/hls"
	"github.com/Lysander66/ace/pkg/logger"
	"github.com/bluenviron/gohlslib/pkg/playlist"
	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/logic"
)

func main() {
	slog.SetDefault(logger.New(slog.LevelInfo))

	srv := logic.NewLalServer()

	rawURL := "http://devimages.apple.com/iphone/samples/bipbop/gear1/prog_index.m3u8"
	go customizePub(srv, "c110", rawURL)

	err := srv.RunLoop()
	slog.Info("Run", "err", err)
}

func customizePub(srv logic.ILalServer, streamName, rawURL string) {
	pubSession, err := srv.AddCustomizePubSession(streamName)
	if err != nil {
		slog.Error("AddCustomizePubSession", "err", err)
		return
	}
	pubSession.WithOption(func(option *base.AvPacketStreamOption) {
		option.VideoFormat = base.AvPacketStreamVideoFormatAnnexb
	})
	pub := hls.NewPublisher(context.Background(), pubSession)

	pullSession := hls.NewPullSession(rawURL, nil, func(_ *playlist.MediaSegment, buf []byte) {
		pub.Write(buf)
	})
	err = pullSession.Start()
	if err != nil {
		slog.Error("customizePub", "err", err)
		return
	}

	slog.Info("customizePub", "err", <-pullSession.Wait())
}
