package crpc

import (
	"context"
	"fmt"

	"github.com/Lysander66/ace/pkg/crpc/sd/etcdv3"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Dial(ctx context.Context, serviceName string, etcdUrls []string) (*grpc.ClientConn, error) {
	target := fmt.Sprintf("etcd:///%s", etcdv3.KeyPrefix(serviceName))
	cli, err := clientv3.NewFromURLs(etcdUrls)
	if err != nil {
		return nil, err
	}
	etcdResolver, err := resolver.NewBuilder(cli)
	if err != nil {
		return nil, err
	}
	return grpc.DialContext(ctx, target,
		grpc.WithResolvers(etcdResolver),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}
