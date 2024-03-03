package streamurl

import (
	"testing"
	"time"
)

func Test_PublishingAddress(t *testing.T) {
	var (
		name        = "aliyun"
		pushDomain  = "demo.aliyundoc.com"
		pushAuthKey = "key123"
		pullDomain  = "example.aliyundoc.com"
		pullAuthKey = "key456"
		appName     = "live"
		streamName  = "test"
		expireAt    = time.Now().Add(24 * time.Hour).Unix()
	)
	agent, err := New(name, pushDomain, pushAuthKey, pullDomain, pullAuthKey, true)
	if err != nil {
		t.Error(err)
		return
	}

	addr := agent.PublishingAddress(appName, streamName, expireAt)
	t.Log(addr)

	flvPlayUrl := agent.FlvPlayUrl(appName, streamName, expireAt)
	t.Log(flvPlayUrl)

	hlsPlayUrl := agent.HlsPlayUrl(appName, streamName, expireAt)
	t.Log(hlsPlayUrl)
}
