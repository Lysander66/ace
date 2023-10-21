package service

import (
	"context"
	"log/slog"

	"github.com/Lysander66/ace/pkg/crpc/_examples/pb"
)

var _ pb.HelloWorldServiceServer = (*HelloSvc)(nil)

type HelloSvc struct {
	pb.UnimplementedHelloWorldServiceServer
}

func (h HelloSvc) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	slog.Info("SayHello", "name", req.Name)
	return &pb.HelloResponse{Message: "hello " + req.Name}, nil
}
