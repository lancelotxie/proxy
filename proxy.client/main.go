/*
 * @Author: EagleXiang
 * @LastEditors  : EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-05 21:14:44
 * @LastEditTime : 2019-12-19 23:26:26
 */

package main

import (
	"context"
	"flag"
	"log"
	"net"
	_ "net/http/pprof"

	"github.com/lancelot/proxy/proxy.client/door"
	"github.com/lancelot/proxy/proxy.client/kcp"
	"github.com/lancelot/proxy/proxy.lib/configuration"
	config "github.com/lancelot/proxy/proxy.lib/configuration/grpc"
	"github.com/lancelot/proxy/proxy.lib/debug"
	"github.com/lancelot/proxy/proxy.lib/path"
	client "github.com/lancelot/proxy/proxy.server/grpc.client"

	_ "github.com/lancelot/proxy/proxy.client/parsers/http"
	_ "github.com/lancelot/proxy/proxy.client/parsers/https"
)

func main() {
	ctx := context.Background()

	dbg := flag.String("debug", "off", "是否开启调试模式")
	port := flag.Int("ctrl-port", 8085, "port for controller")
	flag.Parse()

	// 初始化系统
	initSystem(*dbg)
	defer closeSystem()

	// 开启 KCP 转发服务
	err := kcp.Start(ctx)
	if err != nil {
		log.Println(err)
	}

	// 开启本地代理服务
	err = startProxy()
	if err != nil {
		log.Println(err)
	}

	err = config.Start("127.0.0.1", *port)
	if err != nil {
		log.Println(err)
	}
}

func initSystem(dbg string) {
	debug.TryRun(dbg)

	path.Init(path.Client)

	err := configuration.Load()
	if err != nil {
		log.Println(err)
	}

	err = initGRPCClient()
	if err != nil {
		log.Println(err)
	}
}

func closeSystem() {
	config.Stop()
}

func initGRPCClient() (err error) {
	log.Println("gRPC 目的地址被初始化为：127.0.0.1:9529")
	err = client.Init("127.0.0.1", 9529)
	return
}

// startProxy 开启本地代理服务
func startProxy() (err error) {
	lis, err := net.Listen("tcp", "0.0.0.0:9527")
	if err != nil {
		return
	}
	log.Println("本地代理服务成功监听：", "0.0.0.0:9527")

	d := door.NewDoor(lis)
	go d.Serve()
	return
}
