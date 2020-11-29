package service

import (
	"net"

	mynet "github.com/lancelotXie/proxy/proxy.lib/net"
	dial "github.com/lancelotXie/proxy/proxy.server/dial"
	flow "github.com/lancelotXie/proxy/proxy.server/flow"
	grpc "github.com/lancelotXie/proxy/proxy.server/grpc.server"
)

// NewService : 包外可以使用的 服务器
func NewService(lis net.Listener) mynet.Service {
	baseLis := grpc.NewListener(lis)

	server := &server{
		Dnet:   &dial.WebServer{},
		Tunnel: &flow.Conn2Conn{},
		base:   baseLis,
	}
	return server
}
