package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Lysander66/ace/pkg/common/cnet"
	"github.com/Lysander66/ace/pkg/hls"
	"github.com/Lysander66/ace/pkg/logger"
	"github.com/bluenviron/gohlslib/pkg/playlist"
	"github.com/cheggaaa/pb/v3"
)

var (
	iFlag = flag.String("i", "", "m3u8下载地址")
	oFlag = flag.String("o", "", "自定义文件名不带后缀")
	nFlag = flag.Int("n", 5, "下载线程数(默认5)")
)

func main() {
	slog.SetDefault(logger.New(slog.LevelInfo))
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	dir := *oFlag
	if dir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			slog.Error("UserHomeDir", "err", err)
			return
		}
		slog.Info("UserHomeDir", "HOME", homeDir)
		dir = filepath.Join(filepath.Join(homeDir, "Downloads"), time.Now().Format(time.DateOnly))
	}

	download(*iFlag, dir, *nFlag)
}

func download(rawURL, dir string, n int) {
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		slog.Error("Mkdir", "dir", dir)
		return
	}

	var bar *pb.ProgressBar

	client := cnet.New()
	pullSession := hls.NewPullSession(rawURL, client, func(seg *playlist.MediaSegment, buf []byte) {
		name := seg.URI
		if index := strings.Index(seg.URI, "?"); index != -1 {
			name = seg.URI[0:index]
		}
		if index := strings.LastIndex(name, "/"); index != -1 {
			name = name[index+1:]
		}
		if err := os.WriteFile(filepath.Join(dir, name), buf, os.ModePerm); err != nil {
			slog.Error("WriteFile", "err", err, "name", name)
			return
		}
		bar.Increment()
	})
	pullSession.SetParallel(n)
	pullSession.AfterFirstPlaylistDownload = func(p playlist.Playlist, data []byte) {
		var name string
		switch pl := p.(type) {
		case *playlist.Media:
			bar = pb.StartNew(len(pl.Segments))
			name = "index.m3u8"
		case *playlist.Multivariant:
			name = "variant.m3u8"
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, os.ModePerm); err != nil {
			slog.Error("WriteFile", "err", err, "name", name)
			return
		}
		slog.Info("download", "name", name)
	}

	pullSession.Start()

	err := <-pullSession.Wait()
	slog.Info("done", "err", err, "url", rawURL)

	bar.Finish()
}
