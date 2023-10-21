package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"google.golang.org/grpc"
)

type Option func(o *Server)

func Network(network string) Option {
	return func(s *Server) {
		s.network = network
	}
}

func Port(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// Server is a gRPC server wrapper.
type Server struct {
	*grpc.Server
	baseCtx context.Context
	lis     net.Listener
	err     error
	network string
	port    int
	ready   chan struct{}
}

func NewServer(opts ...Option) *Server {
	srv := &Server{
		baseCtx: context.Background(),
		network: "tcp",
		ready:   make(chan struct{}),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.Server = grpc.NewServer()
	return srv
}

func (s *Server) Start() error {
	lis, err := net.Listen(s.network, ":"+strconv.Itoa(s.port))
	if err != nil {
		return err
	}

	if s.port == 0 {
		addr := lis.Addr().(*net.TCPAddr)
		s.port = addr.Port
	}
	close(s.ready)

	slog.Info(fmt.Sprintf("Listen %s :%d", s.network, s.port))
	return s.Serve(lis)
}

func (s *Server) Stop() error {
	s.GracefulStop()
	slog.Info("[gRPC] server stopping")
	return nil
}

func (s *Server) Port() int {
	// wait for net.Listen to automatically choose a port number.
	if s.port == 0 {
		<-s.ready
	}
	return s.port
}
