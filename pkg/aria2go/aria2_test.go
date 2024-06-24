package aria2go

import (
	"log"
	"testing"
)

var testClient, _ = NewClient("http://localhost:6800/jsonrpc", "", nil)

type DummyNotifier struct{}

func (DummyNotifier) OnDownloadStart(events []Event)      { log.Printf("%s started.", events) }
func (DummyNotifier) OnDownloadPause(events []Event)      { log.Printf("%s paused.", events) }
func (DummyNotifier) OnDownloadStop(events []Event)       { log.Printf("%s stopped.", events) }
func (DummyNotifier) OnDownloadComplete(events []Event)   { log.Printf("%s completed.", events) }
func (DummyNotifier) OnDownloadError(events []Event)      { log.Printf("%s error.", events) }
func (DummyNotifier) OnBtDownloadComplete(events []Event) { log.Printf("bt %s completed.", events) }

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

	client, _ := NewClient("http://localhost:6800/jsonrpc", "rpcSecret", DummyNotifier{})
	gid, err := client.AddURI(uris, options)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("gid", gid)

	select {}
}

func TestClient_GetGlobalStat(t *testing.T) {
	globalStat, err := testClient.GetGlobalStat()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v\n", globalStat)
}

func TestClient_ListMethods(t *testing.T) {
	methods, err := testClient.ListMethods()
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(methods)
	t.Log(len(methods))
}
