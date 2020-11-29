package kcp

import (
	"context"
	"log"

	"github.com/lancelotXie/proxy/proxy.lib/net/kcp"
	"github.com/lancelotXie/proxy/proxy.lib/net/relay"
)

// Start 开始 kcp 转发服务
func Start(ctx context.Context, localAddr string) (err error) {
	lis, err := kcp.Listen("0.0.0.0:9588")
	if err != nil {
		return
	}
	log.Println("KCP 转发服务成功监听：0.0.0.0:9588")

	r := relay.NewTCP(ctx, localAddr, lis)
	err = r.Start()
	return
}
