package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdClient 封装 etcd 客户端配置和连接
type EtcdClient struct {
	client *clientv3.Client
}

// ClientConfig 定义 etcd 客户端配置
type ClientConfig struct {
	Endpoints   []string
	Username    string
	Password    string
	DialTimeout time.Duration
	TLS         *tls.Config
}

// SyncStats 同步统计信息
type SyncStats struct {
	TotalKeys    int
	SuccessCount int
	FailCount    int
}

// NewEtcdClient 创建新的 etcd 客户端
func NewEtcdClient(cfg ClientConfig) (*EtcdClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		Username:    cfg.Username,
		Password:    cfg.Password,
		DialTimeout: cfg.DialTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &EtcdClient{client: client}, nil
}

// WithTLS 添加 TLS 配置
func (c *ClientConfig) WithTLS(certFile, keyFile, ca string) error {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return fmt.Errorf("failed to load cert: %w", err)
	}

	caCert, err := os.ReadFile(ca)
	if err != nil {
		return fmt.Errorf("failed to read ca cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	c.TLS = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	return nil
}

// Close 关闭客户端连接
func (e *EtcdClient) Close() error {
	return e.client.Close()
}

// SyncKeysByPrefix 同步指定前缀的键值对
func SyncKeysByPrefix(ctx context.Context, src, dst *EtcdClient, prefix string) (*SyncStats, error) {
	stats := &SyncStats{}

	// 获取源 etcd 中指定前缀的所有键值对
	resp, err := src.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from source: %w", err)
	}

	stats.TotalKeys = len(resp.Kvs)

	// 同步到目标 etcd
	for _, kv := range resp.Kvs {
		if err := syncSingleKey(ctx, dst.client, kv); err != nil {
			log.Printf("Failed to sync key %s: %v", string(kv.Key), err)
			stats.FailCount++
			continue
		}
		stats.SuccessCount++
	}

	return stats, nil
}

// syncSingleKey 同步单个键值对
func syncSingleKey(ctx context.Context, dstClient *clientv3.Client, kv *mvccpb.KeyValue) error {
	_, err := dstClient.Put(ctx, string(kv.Key), string(kv.Value))
	return err
}

func main() {
	// 命令行参数定义
	srcEndpoint := flag.String("src", "localhost:2379", "源 etcd 地址")
	dstEndpoint := flag.String("dst", "localhost:2379", "目标 etcd 地址")
	prefix := flag.String("prefix", "/", "要同步的键前缀")
	timeout := flag.Duration("timeout", 5*time.Second, "操作超时时间")
	flag.Parse()

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// 初始化源客户端
	srcClient, err := NewEtcdClient(ClientConfig{
		Endpoints:   []string{*srcEndpoint},
		DialTimeout: *timeout,
	})
	if err != nil {
		log.Fatalf("Failed to create source client: %v", err)
	}
	defer srcClient.Close()

	// 初始化目标客户端
	dstClient, err := NewEtcdClient(ClientConfig{
		Endpoints:   []string{*dstEndpoint},
		DialTimeout: *timeout,
	})
	if err != nil {
		log.Fatalf("Failed to create destination client: %v", err)
	}
	defer dstClient.Close()

	// 执行同步
	stats, err := SyncKeysByPrefix(ctx, srcClient, dstClient, *prefix)
	if err != nil {
		log.Fatalf("Sync failed: %v", err)
	}

	// 打印结果
	fmt.Printf("Sync completed:\n")
	fmt.Printf("Total keys: %d\n", stats.TotalKeys)
	fmt.Printf("Successfully synced: %d\n", stats.SuccessCount)
	fmt.Printf("Failed: %d\n", stats.FailCount)
}
