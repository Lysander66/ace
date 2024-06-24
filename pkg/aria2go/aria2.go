package aria2go

import (
	"github.com/Lysander66/ace/pkg/jsonrpc"
)

const (
	methodAddUri               = "aria2.addUri"
	methodAddTorrent           = "aria2.addTorrent"
	methodGetPeers             = "aria2.getPeers"
	methodAddMetalink          = "aria2.addMetalink"
	methodRemove               = "aria2.remove"
	methodPause                = "aria2.pause"
	methodForcePause           = "aria2.forcePause"
	methodPauseAll             = "aria2.pauseAll"
	methodForcePauseAll        = "aria2.forcePauseAll"
	methodUnpause              = "aria2.unpause"
	methodUnpauseAll           = "aria2.unpauseAll"
	methodForceRemove          = "aria2.forceRemove"
	methodChangePosition       = "aria2.changePosition"
	methodTellStatus           = "aria2.tellStatus"
	methodGetUris              = "aria2.getUris"
	methodGetFiles             = "aria2.getFiles"
	methodGetServers           = "aria2.getServers"
	methodTellActive           = "aria2.tellActive"
	methodTellWaiting          = "aria2.tellWaiting"
	methodTellStopped          = "aria2.tellStopped"
	methodGetOption            = "aria2.getOption"
	methodChangeUri            = "aria2.changeUri"
	methodChangeOption         = "aria2.changeOption"
	methodGetGlobalOption      = "aria2.getGlobalOption"
	methodChangeGlobalOption   = "aria2.changeGlobalOption"
	methodPurgeDownloadResult  = "aria2.purgeDownloadResult"
	methodRemoveDownloadResult = "aria2.removeDownloadResult"
	methodGetVersion           = "aria2.getVersion"
	methodGetSessionInfo       = "aria2.getSessionInfo"
	methodShutdown             = "aria2.shutdown"
	methodForceShutdown        = "aria2.forceShutdown"
	methodGetGlobalStat        = "aria2.getGlobalStat"
	methodSaveSession          = "aria2.saveSession"
	methodMultiCall            = "system.multicall"
	methodListMethods          = "system.listMethods"
	methodListNotifications    = "system.listNotifications"
)

type Client struct {
	secret    string
	rpcClient *jsonrpc.Client
}

type Option func(o *Client)

func NewClient(endpoint, rpcSecret string) *Client {
	c := &Client{
		secret:    rpcSecret,
		rpcClient: jsonrpc.NewClient(endpoint),
	}
	return c
}

func (c *Client) token() string {
	return "token:" + c.secret
}

// AddURI
// https://aria2.github.io/manual/en/html/aria2c.html#aria2.addUri
func (c *Client) AddURI(id string, uris []string, options ...any) (string, error) {
	var params []any
	if c.secret != "" {
		params = append(params, c.token())
	}
	params = append(params, uris)
	if options != nil {
		params = append(params, options...)
	}

	req := jsonrpc.NewRequest(methodAddUri, params, id)
	resp, err := c.rpcClient.Call(req)
	if err != nil {
		return "", err
	}

	return resp.GetString()
}

func (c *Client) ListMethods(id string) (methods []string, err error) {
	req := jsonrpc.NewRequest(methodListMethods, nil, id)
	resp, err := c.rpcClient.Call(req)
	if err != nil {
		return nil, err
	}

	err = resp.GetAny(&methods)
	return
}
