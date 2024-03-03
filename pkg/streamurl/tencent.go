package streamurl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Lysander66/ace/pkg/common"
)

// TxGeneratort 腾讯
// https://cloud.tencent.com/document/product/267/35257
type TxGeneratort struct{}

func (t TxGeneratort) PublishingAddress(key, app, stream string, exp int64) string {
	return fmt.Sprintf("/%s/%s", app, stream) + t.secret(key, stream, exp)
}

func (t TxGeneratort) FlvPlayUrl(key, app, stream string, exp int64) string {
	return fmt.Sprintf("/%s/%s.%s", app, stream, "flv") + t.secret(key, stream, exp)
}

func (t TxGeneratort) HlsPlayUrl(key, app, stream string, exp int64) string {
	return fmt.Sprintf("/%s/%s.%s", app, stream, "m3u8") + t.secret(key, stream, exp)
}

func (t TxGeneratort) secret(key, stream string, exp int64) string {
	if key != "" && exp > 0 {
		txTime := strings.ToUpper(strconv.FormatInt(exp, 16))
		txSecret := common.MD5Sum(key + stream + txTime)
		return "?txSecret=" + txSecret + "&txTime=" + txTime
	}
	return ""
}
