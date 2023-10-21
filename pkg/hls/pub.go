package hls

import (
	"bytes"
	"context"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/logic"
	"github.com/yapingcat/gomedia/go-codec"
	"github.com/yapingcat/gomedia/go-mpeg2"
)

type Publisher struct {
	ctx     context.Context
	ss      logic.ICustomizePubSessionContext
	demuxer *mpeg2.TSDemuxer
}

func NewPublisher(ctx context.Context, session logic.ICustomizePubSessionContext) *Publisher {
	foundAudio := false
	demuxer := mpeg2.NewTSDemuxer()
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		packet := base.AvPacket{
			Timestamp: int64(dts),
			Pts:       int64(pts),
		}
		switch cid {
		case mpeg2.TS_STREAM_AAC:
			if !foundAudio {
				foundAudio = true
				asc, _ := codec.ConvertADTSToASC(frame)
				session.FeedAudioSpecificConfig(asc.Encode())
			}

			packet.PayloadType = base.AvPacketPtAac
			packet.Payload = frame[7:]
			session.FeedAvPacket(packet)

		case mpeg2.TS_STREAM_H264:
			packet.PayloadType = base.AvPacketPtAvc
			packet.Payload = frame
			session.FeedAvPacket(packet)

		case mpeg2.TS_STREAM_H265:
			packet.PayloadType = base.AvPacketPtHevc
			packet.Payload = frame
			session.FeedAvPacket(packet)

		case mpeg2.TS_STREAM_AUDIO_MPEG1, mpeg2.TS_STREAM_AUDIO_MPEG2:
			//if !foundAudio { foundAudio = true }
		}
	}

	pub := &Publisher{
		ctx:     ctx,
		ss:      session,
		demuxer: demuxer,
	}
	return pub
}

func (p *Publisher) Write(msg []byte) error {
	return p.demuxer.Input(bytes.NewReader(msg))
}
