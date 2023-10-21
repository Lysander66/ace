package ws

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
)

const reconnectInterval = 2 * time.Second

type MsgHandleFunc func(msg []byte)

type SimpleWsClientParam struct {
	Url        string
	HandleFunc MsgHandleFunc
}

type SimpleWsClient struct {
	token     string
	reconnect chan SimpleWsClientParam
}

func NewSimpleWsClient(token string) *SimpleWsClient {
	s := &SimpleWsClient{
		token:     token,
		reconnect: make(chan SimpleWsClientParam, 1),
	}
	go func() {
		for {
			select {
			case param := <-s.reconnect:
				slog.Info("reconnect", "url", param.Url)
				go s.connect(param)
			}
		}
	}()
	return s
}

func (s SimpleWsClient) Subscribe(params ...SimpleWsClientParam) {
	for _, param := range params {
		go s.connect(param)
	}
}

func (s SimpleWsClient) connect(param SimpleWsClientParam) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	url := param.Url
	c, _, err := websocket.DefaultDialer.Dial(url, http.Header{"token": []string{s.token}})
	if err != nil {
		slog.Error("connect fail.", "err", err, "url", url)
		s.reconnect <- param
		return
	}
	defer c.Close()
	slog.Info("connect success!", "url", url)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				slog.Error("disconnect", "err", err, "url", url)

				time.Sleep(reconnectInterval)
				s.reconnect <- param
				return
			}
			param.HandleFunc(message)
		}
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			err = c.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				slog.Error("WriteMessage", "err", err)
				return
			}
		case <-interrupt:
			err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				slog.Error("interrupt", "err", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
