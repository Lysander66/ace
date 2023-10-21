package etcdv3

import (
	"log/slog"
	"sync"

	"github.com/Lysander66/ace/pkg/crpc/sd"
)

type Registrar struct {
	client  Client
	quitmtx sync.Mutex
	quit    chan struct{}
}

func NewRegistrar(client Client) *Registrar {
	return &Registrar{client: client}
}

func (r *Registrar) Register(s *sd.Service) {
	if err := r.client.Register(s); err != nil {
		slog.Error("Register", "err", err)
		return
	}
	if s.TTL != nil {
		slog.Info("Register", "lease", r.client.LeaseID())
	} else {
		slog.Info("Register")
	}
}

func (r *Registrar) Deregister(s *sd.Service) {
	if err := r.client.Deregister(s); err != nil {
		slog.Error("Deregister", "err", err)
	} else {
		slog.Info("Deregister")
	}

	r.quitmtx.Lock()
	defer r.quitmtx.Unlock()
	if r.quit != nil {
		close(r.quit)
		r.quit = nil
	}
}

func KeyPrefix(key string) string {
	return "services/" + key
}
