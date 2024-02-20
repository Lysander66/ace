package streamurl

import (
	"fmt"
	"strings"
)

// Generator StreamURLGenerator
type Generator interface {
	PublishingAddress(key, app, stream string, exp int64) string //推流地址
	FlvPlayUrl(key, app, stream string, exp int64) string        //播放地址
	HlsPlayUrl(key, app, stream string, exp int64) string        //播放地址
}

type Agent struct {
	Name        string
	pushDomain  string
	pushAuthKey string
	pullDomain  string
	pullAuthKey string
	https       bool
	generator   Generator
}

func New(name, pushDomain, pushAuthKey, pullDomain, pullAuthKey string, https bool) (*Agent, error) {
	agent := &Agent{
		Name:        strings.TrimSpace(name),
		pushDomain:  strings.TrimSpace(pushDomain),
		pushAuthKey: strings.TrimSpace(pushAuthKey),
		pullDomain:  strings.TrimSpace(pullDomain),
		pullAuthKey: strings.TrimSpace(pullAuthKey),
		https:       https,
	}
	switch agent.Name {
	case "aliyun":
		agent.generator = &AliGeneratort{}
	case "tencent":
		agent.generator = &TxGeneratort{}
	case "wangsu":
		agent.generator = &WsGeneratort{}
	default:
		return nil, fmt.Errorf("%s is not supported yet", name)
	}
	return agent, nil
}

func (agt *Agent) PublishingAddress(appName, streamName string, expireAt int64) string {
	return "rtmp://" + agt.pushDomain + agt.generator.PublishingAddress(agt.pushAuthKey, appName, streamName, expireAt)
}

func (agt *Agent) FlvPlayUrl(appName, streamName string, expireAt int64) string {
	path := agt.pullDomain + agt.generator.FlvPlayUrl(agt.pullAuthKey, appName, streamName, expireAt)
	if agt.https {
		return "https://" + path
	}
	return "http://" + path
}

func (agt *Agent) HlsPlayUrl(appName, streamName string, expireAt int64) string {
	path := agt.pullDomain + agt.generator.HlsPlayUrl(agt.pullAuthKey, appName, streamName, expireAt)
	if agt.https {
		return "https://" + path
	}
	return "http://" + path
}
