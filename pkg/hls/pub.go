package hls

import (
	"bytes"
	"context"

	"github.com/q191201771/lal/pkg/aac"
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
	audioSampleRate := 0
	demuxer := mpeg2.NewTSDemuxer()
	demuxer.OnFrame = func(cid mpeg2.TS_STREAM_TYPE, frame []byte, pts uint64, dts uint64) {
		switch cid {
		case mpeg2.TS_STREAM_AAC:
			if !foundAudio {
				asc, err := codec.ConvertADTSToASC(frame)
				if err != nil {
					return
				}
				session.FeedAudioSpecificConfig(asc.Encode())
				audioSampleRate = codec.AACSampleIdxToSample(int(asc.Sample_freq_index))
				foundAudio = true
			}

			var preAudioDts int64
			ctx := aac.AdtsHeaderContext{}
			for len(frame) > aac.AdtsHeaderLength {
				ctx.Unpack(frame[:])
				if preAudioDts == 0 {
					preAudioDts = int64(dts)
				} else {
					preAudioDts += 1024 * 1000 / int64(audioSampleRate)
				}

				aacPacket := base.AvPacket{
					Timestamp:   preAudioDts,
					PayloadType: base.AvPacketPtAac,
				}
				payload := frame[aac.AdtsHeaderLength:ctx.AdtsLength]
				frame = frame[ctx.AdtsLength:]
				aacPacket.Payload = payload
				session.FeedAvPacket(aacPacket)
			}

		case mpeg2.TS_STREAM_H264:
			packet := base.AvPacket{
				PayloadType: base.AvPacketPtAvc,
				Timestamp:   int64(dts),
				Pts:         int64(pts),
				Payload:     frame,
			}
			session.FeedAvPacket(packet)

		case mpeg2.TS_STREAM_H265:
			packet := base.AvPacket{
				PayloadType: base.AvPacketPtHevc,
				Timestamp:   int64(dts),
				Pts:         int64(pts),
				Payload:     frame,
			}
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
