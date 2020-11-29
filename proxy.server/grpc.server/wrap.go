package server

import (
	"context"
	"net"

	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
)

// Service : GRPC服务器对外接口
type Service interface {
	Addr() net.Addr
	Close() error
	Accept() (context.Context, net.Conn, error)
}

// NewListener : 构造新grpc服务。方法内部有监听，直到报错才会返回
func NewListener(lis net.Listener) Service {
	gserver := newgrpcServer(lis)
	proto.RegisterStreamServiceServer(gserver.server, gserver)
	go gserver.start()
	return gserver
}
