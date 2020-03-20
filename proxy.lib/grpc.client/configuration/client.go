package configuration

import (
	"context"
	"fmt"

	rpcconfig "github.com/lancelot/proxy/proxy.lib/configuration/grpc/config"

	"google.golang.org/grpc"
)

var defaultClient rpcconfig.ConfigurationServerClient

// Init 初始化客 gRPC 户端
func Init(ip string, port int) (err error) {
	host := fmt.Sprint(ip, ":", port)
	cc, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return
	}

	defaultClient = rpcconfig.NewConfigurationServerClient(cc)
	return
}

// Set 设置 key/value
func Set(ctx context.Context, key string, value interface{}) (err error) {
	kv := rpcconfig.KeyValue{
		Key:   key,
		Value: value,
	}
	kvs := rpcconfig.KeyValues{}.Push(kv)

	in := new(rpcconfig.Content)
	in.Content = kvs.Bytes()
	_, err = defaultClient.Set(ctx, in)
	return
}

// Get 获取 value
func Get(ctx context.Context, key string) (value interface{}, err error) {
	kv := rpcconfig.KeyValue{
		Key: key,
	}
	kvs := rpcconfig.KeyValues{}.Push(kv)

	in := new(rpcconfig.Content)
	in.Content = kvs.Bytes()
	r, err := defaultClient.Get(ctx, in)

	content := r.GetContent()
	kvs, err = rpcconfig.ParseKeyValues(content)
	if err != nil {
		return
	}

	if len(kvs) > 0 {
		kv = kvs[0]
		value = kv.Value
	}
	return
}

// Save 保存配置
func Save(ctx context.Context) (err error) {
	in := new(rpcconfig.Nop)
	_, err = defaultClient.Save(ctx, in)
	return
}
