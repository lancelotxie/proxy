package kcp

import (
	"context"
	"log"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/configuration"
	"github.com/lancelotXie/proxy/proxy.lib/net/kcp"
	"github.com/lancelotXie/proxy/proxy.lib/net/relay"

	"github.com/pkg/errors"
)

// Start 开始 KCP 转发服务
func Start(ctx context.Context) (err error) {
	addr, ok := configuration.GetString("remote.addr")
	if !ok {
		log.Println("缺失 KCP 远端地址")
		return
	}
	log.Println("KCP 远端地址为：", addr)

	lis, err := net.Listen("tcp", "127.0.0.1:9529")
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	log.Println("KCP 转发服务成功监听：127.0.0.1:9529")

	dialer := func() (conn net.Conn, err error) {
		addr, ok := configuration.GetString("remote.addr")
		if !ok {
			err = errors.WithStack(configuration.ErrRemoteAddrNotFound)
			return
		}

		conn, err = kcp.Dial(addr)
		return
	}

	s := relay.New(ctx, dialer, lis)
	s.Start()
	return
}
