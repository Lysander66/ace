package streamurl

import (
	"fmt"
	"strconv"

	"github.com/Lysander66/ace/pkg/cryptogo"
)

// WsGeneratort 网宿
// https://www.wangsu.com/document/livestream/171
type WsGeneratort struct{}

func (w WsGeneratort) PublishingAddress(key, app, stream string, exp int64) string {
	return fmt.Sprintf("/%s/%s", app, stream) + w.secret(key, app, stream, exp)
}

func (w WsGeneratort) FlvPlayUrl(key, app, stream string, exp int64) string {
	return fmt.Sprintf("/%s/%s.%s", app, stream, "flv") + w.secret(key, app, stream, exp)
}

func (w WsGeneratort) HlsPlayUrl(key, app, stream string, exp int64) string {
	return fmt.Sprintf("/%s/%s.%s", app, stream, "m3u8") + w.secret(key, app, stream, exp)
}

func (w WsGeneratort) secret(key, app, stream string, exp int64) string {
	if key != "" && exp > 0 {
		wsTime := strconv.FormatInt(exp, 16)
		wsSecret := cryptogo.MD5Sum(fmt.Sprintf("/%s/%s%s%s", app, stream, key, wsTime))
		return "?wsSecret=" + wsSecret + "&wsABSTime=" + wsTime
	}
	return ""
}
