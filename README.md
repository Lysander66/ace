# Ace

Using etcd discovery with go-grpc

[B 站视频: gRPC + etcd 服务注册与发现](https://www.bilibili.com/video/BV1sP411a7B9/)

## Getting Started

addsvc.go

```go
func main() {
	svc := ace.NewService(
		ace.Name("hello"),
		ace.EtcdUrls([]string{"localhost:2379"}),
		ace.RegisterServiceFn(func(s *server.Server) {
			pb.RegisterAddServer(s, &service.AddSvc{})
			pb.RegisterHelloWorldServiceServer(s, &service.HelloSvc{})
		}),
	)

	svc.Run()
}
```

addcli.go

```go
func main() {
	ctx := context.Background()
	conn, err := ace.Dial(ctx, "hello", []string{"localhost:2379"})
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	var a, b int64 = 2, 3
	addClient := pb.NewAddClient(conn)
	sumReply, err := addClient.Sum(ctx, &pb.SumRequest{A: a, B: b})
	if err != nil {
		fmt.Println("err", err)
		return
	}
	fmt.Printf("%d + %d = %d\n", a, b, sumReply.V)
}
```

## etcd 与 gRPC 版本兼容问题

> go.etcd.io/etcd/client/v3@v3.5.9/naming/resolver/resolver.go:22:11:
> cannot use target.Endpoint (value of type func() string) as type string in struct literal

- [cannot use target.Endpoint as type string in struct literal](https://github.com/etcd-io/etcd/issues/15286)

```
replace google.golang.org/grpc => google.golang.org/grpc v1.52.3
```

## references

1. [gRPC naming and discovery](https://etcd.io/docs/v3.5/dev-guide/grpc_naming)
1. [go-kit](https://github.com/go-kit/kit/tree/master/sd/etcdv3)
1. [go-micro](https://github.com/go-micro/go-micro)
