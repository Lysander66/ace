package crpc

import (
	"context"

	"github.com/Lysander66/ace/pkg/crpc/sd"
	"github.com/Lysander66/ace/pkg/crpc/server"
)

type options struct {
	name              string
	metadata          map[string]string
	ctx               context.Context
	server            *server.Server
	etcdUrls          []string
	registrar         sd.Registrar
	registerServiceFn func(*server.Server)
	serviceInstance   *sd.Service
	getMachineIPFn    GetMachineIP
	beforeStart       []func() error
	afterStart        []func() error
}

type GetMachineIP func() string

type Option func(o *options)

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func Server(srv *server.Server) Option {
	return func(o *options) { o.server = srv }
}

func EtcdUrls(urls []string) Option {
	return func(o *options) { o.etcdUrls = urls }
}

func Registrar(r sd.Registrar) Option {
	return func(o *options) { o.registrar = r }
}

func RegisterServiceFn(fn func(*server.Server)) Option {
	return func(o *options) {
		o.registerServiceFn = fn
	}
}

func GetMachineIPFn(fn GetMachineIP) Option {
	return func(o *options) {
		o.getMachineIPFn = fn
	}
}

func BeforeStart(fn func() error) Option {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

func AfterStart(fn func() error) Option {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}
