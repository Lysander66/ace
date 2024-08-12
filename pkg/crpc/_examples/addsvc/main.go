package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Lysander66/ace/pkg/crpc"
	"github.com/Lysander66/ace/pkg/crpc/_examples/addsvc/service"
	"github.com/Lysander66/ace/pkg/crpc/_examples/pb"
	"github.com/Lysander66/ace/pkg/crpc/server"
	"github.com/Lysander66/zephyr/pkg/logger"
)

const HelloService = "hello"

func main() {
	slog.SetDefault(logger.New(slog.LevelInfo))

	etcdUrls := []string{"localhost:2379"}
	if s := os.Getenv("ACE_REGISTRY_ADDRESS"); s != "" {
		etcdUrls = strings.Split(s, ",")
	}

	svc := crpc.NewService(
		crpc.Name(HelloService),
		//crpc.Server(server.NewServer(server.Port(8081))),
		crpc.EtcdUrls(etcdUrls),
		crpc.RegisterServiceFn(func(s *server.Server) {
			pb.RegisterAddServer(s, &service.AddSvc{})
			pb.RegisterHelloWorldServiceServer(s, &service.HelloSvc{})
		}),
	)

	svc.Run()
}
