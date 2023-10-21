package sd

import "time"

const minHeartBeatTime = 500 * time.Millisecond

// Registrar registers instance information to a service discovery system when
// an instance becomes alive and healthy, and deregisters that information when
// the service becomes unhealthy or goes away.
//
// Registrar implementations exist for various service discovery systems. Note
// that identifying instance information (e.g. host:port) must be given via the
// concrete constructor; this interface merely signals lifecycle changes.
type Registrar interface {
	Register(*Service)
	Deregister(*Service)
}

// Service holds the instance identifying data you want to publish to etcd. Key
// must be unique, and value is the string returned to subscribers, typically
// called the "instance" string in other parts of package sd.
type Service struct {
	Key   string // unique key, e.g. "/service/foobar/1.2.3.4:8080"
	Value string // returned to subscribers, e.g. "http://1.2.3.4:8080"
	TTL   *TTLOption
}

func NewService(key, value string) *Service {
	return &Service{
		Key:   key,
		Value: value,
		TTL:   NewTTLOption(3*time.Second, 10*time.Second),
	}
}

// TTLOption allow setting a key with a TTL. This option will be used by a loop
// goroutine which regularly refreshes the lease of the key.
type TTLOption struct {
	heartbeat time.Duration // e.g. time.Second * 3
	TTL       time.Duration // e.g. time.Second * 10
}

// NewTTLOption returns a TTLOption that contains proper TTL settings. Heartbeat
// is used to refresh the lease of the key periodically; its value should be at
// least 500ms. TTL defines the lease of the key; its value should be
// significantly greater than heartbeat.
//
// Good default values might be 3s heartbeat, 10s TTL.
func NewTTLOption(heartbeat, ttl time.Duration) *TTLOption {
	if heartbeat <= minHeartBeatTime {
		heartbeat = minHeartBeatTime
	}
	if ttl <= heartbeat {
		ttl = 3 * heartbeat
	}
	return &TTLOption{
		heartbeat: heartbeat,
		TTL:       ttl,
	}
}
