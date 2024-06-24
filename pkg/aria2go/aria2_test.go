package aria2go

import (
	"testing"
)

const (
	endpoint  = "http://localhost:6800/jsonrpc"
	rpcSecret = ""
	id        = "test"
)

func TestClient_AddURI(t *testing.T) {
	var (
		uris = []string{
			"https://github.com/prometheus/prometheus/releases/download/v2.53.0/prometheus-2.53.0.darwin-amd64.tar.gz",
		}
		options = map[string]any{
			"dir":                       "/root/downloads",
			"out":                       "prometheus-darwin.tar.gz",
			"http-proxy":                "http://127.0.0.1:7890",
			"https-proxy":               "http://127.0.0.1:7890",
			"split":                     5,
			"max-connection-per-server": 1,
		}
	)

	client := NewClient(endpoint, rpcSecret)
	gid, err := client.AddURI(id, uris, options)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(gid)
}

func TestClient_ListMethods(t *testing.T) {
	client := NewClient(endpoint, rpcSecret)
	methods, err := client.ListMethods(id)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(methods)
	t.Log(len(methods))
}
