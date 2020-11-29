package kcp

import (
	"context"
	"log"

	"github.com/lancelotXie/proxy/proxy.lib/net/kcp"
	"github.com/lancelotXie/proxy/proxy.lib/net/relay"
)

// Start 开始 kcp 转发服务
func Start(ctx context.Context, localAddr string, remoteAddr string) (err error) {
	lis, err := kcp.Listen(localAddr)
	if err != nil {
		return
	}
	log.Println("KCP 转发服务成功监听：", localAddr)

	r := relay.NewTCP(ctx, remoteAddr, lis)
	err = r.Start()
	return
}
