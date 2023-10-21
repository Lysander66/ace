package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Lysander66/ace/pkg/crpc"
	"github.com/Lysander66/ace/pkg/crpc/_examples/pb"
	"github.com/Lysander66/ace/pkg/logger"
)

const HelloService = "hello"

var (
	addClient   pb.AddClient
	helloClient pb.HelloWorldServiceClient
)

func main() {
	slog.SetDefault(logger.New(slog.LevelInfo))

	etcdUrls := []string{"localhost:2379"}
	if s := os.Getenv("ACE_REGISTRY_ADDRESS"); s != "" {
		etcdUrls = strings.Split(s, ",")
	}

	ctx := context.Background()

	conn, err := crpc.Dial(ctx, HelloService, etcdUrls)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	addClient = pb.NewAddClient(conn)
	helloClient = pb.NewHelloWorldServiceClient(conn)

	helloResp, err := helloClient.SayHello(ctx, &pb.HelloRequest{Name: "Shakespeare"})
	if err != nil {
		slog.Error("SayHello", "err", err)
		return
	}
	slog.Info(helloResp.Message)

	var a int64 = 2
	var b int64 = 3
	sumReply, err := addClient.Sum(ctx, &pb.SumRequest{A: a, B: b})
	if err != nil {
		slog.Error("Sum", "err", err)
		return
	}
	slog.Info(fmt.Sprintf("%d + %d = %d", a, b, sumReply.V))
}
