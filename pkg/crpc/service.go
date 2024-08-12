package crpc

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Lysander66/ace/pkg/common"
	"github.com/Lysander66/ace/pkg/crpc/sd"
	"github.com/Lysander66/ace/pkg/crpc/sd/etcdv3"
	"github.com/Lysander66/ace/pkg/crpc/server"
	"github.com/Lysander66/zephyr/pkg/znet"
	"google.golang.org/grpc"
)

type Service struct {
	opts   options
	ctx    context.Context
	cancel func()
}

func (s *Service) Name() string { return s.opts.name }

func (s *Service) Metadata() map[string]string { return s.opts.metadata }

func (s *Service) Server() *server.Server { return s.opts.server }

func NewService(opts ...Option) *Service {
	o := options{
		ctx:            context.Background(),
		getMachineIPFn: znet.IntranetIP,
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.server == nil {
		o.server = server.NewServer()
	}

	client, err := etcdv3.NewClient(o.ctx, o.etcdUrls, etcdv3.ClientOptions{
		DialOptions: []grpc.DialOption{grpc.WithBlock()},
	})
	if err != nil {
		panic(err)
	}
	o.registrar = etcdv3.NewRegistrar(client)

	ctx, cancel := context.WithCancel(o.ctx)
	return &Service{opts: o, ctx: ctx, cancel: cancel}
}

func (s *Service) Run() error {
	if err := s.start(); err != nil {
		return err
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, common.ShutdownSignals()...)

	select {
	// wait on kill signal
	case <-ch:
	// wait on context cancel
	case <-s.opts.ctx.Done():
	}

	return s.stop()
}

func (s *Service) start() error {
	slog.Info(fmt.Sprintf("Starting [service] %s", s.Name()))

	for _, fn := range s.opts.beforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	if s.opts.registerServiceFn != nil {
		s.opts.registerServiceFn(s.opts.server)
	}

	go func() {
		if err := s.opts.server.Start(); err != nil {
			panic(err)
		}
	}()

	s.opts.serviceInstance = sd.NewService(
		etcdv3.KeyPrefix(s.opts.name),
		fmt.Sprintf("%s:%d", s.opts.getMachineIPFn(), s.opts.server.Port()),
	)

	s.opts.registrar.Register(s.opts.serviceInstance)

	for _, fn := range s.opts.afterStart {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) stop() error {
	s.opts.registrar.Deregister(s.opts.serviceInstance)
	if s.cancel != nil {
		s.cancel()
	}
	return s.opts.server.Stop()
}
