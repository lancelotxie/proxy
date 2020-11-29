package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/configuration"
	rpcconfig "github.com/lancelotXie/proxy/proxy.lib/configuration/grpc/config"

	"google.golang.org/grpc"
)

type configurationServer struct{}

func (c *configurationServer) Set(ctx context.Context, in *rpcconfig.Content) (out *rpcconfig.Content, err error) {
	out = new(rpcconfig.Content)

	content := in.GetContent()
	kvs, err := rpcconfig.ParseKeyValues(content)
	if err != nil {
		return
	}

	var rs rpcconfig.KeyValues
	for _, kv := range kvs {
		err = configuration.Set(kv.Key, kv.Value)
		if err != nil {
			break
		}
		rs = rs.Push(kv)
	}

	out.Content = rs.Bytes()

	return
}

func (c *configurationServer) Get(ctx context.Context, in *rpcconfig.Content) (out *rpcconfig.Content, err error) {
	out = new(rpcconfig.Content)

	content := in.GetContent()
	kvs, err := rpcconfig.ParseKeyValues(content)
	if err != nil {
		return
	}

	var rs rpcconfig.KeyValues
	for _, kv := range kvs {
		v, ok := configuration.Get(kv.Key)
		if ok {
			kv.Value = v
			rs = rs.Push(kv)
		}
	}

	out.Content = rs.Bytes()

	return
}

// Save 保存配置
func (c *configurationServer) Save(ctx context.Context, in *rpcconfig.Nop) (out *rpcconfig.Nop, err error) {
	out = new(rpcconfig.Nop)
	err = configuration.Save()
	return
}

// Server 对外暴露的服务
type Server struct {
	s *grpc.Server
}

// New 构造新的 Server
func New() (s *Server) {
	s = new(Server)
	s.s = grpc.NewServer()

	rpcconfig.RegisterConfigurationServerServer(s.s, new(configurationServer))

	return
}

// Start 开始服务
func (s *Server) Start(host string, port int) (err error) {
	host = fmt.Sprintf("%s:%d", host, port)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		return
	}

	err = s.s.Serve(lis)
	return
}

// Stop 停止服务
func (s *Server) Stop() {
	s.s.Stop()
}

var defaultServer *Server

// Start 开始 gRPC 服务
func Start(ip string, port int) (err error) {
	defaultServer = New()
	err = defaultServer.Start(ip, port)
	return
}

// Stop 停止 gRPC 服务
func Stop() {
	defaultServer.Stop()
}
