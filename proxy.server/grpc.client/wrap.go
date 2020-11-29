package client

import (
	"context"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/location"

	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	"github.com/lancelotXie/proxy/proxy.lib/dns"
	locationbase "github.com/lancelotXie/proxy/proxy.lib/location/base"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
)

var defaultPool *pool
var defaultDNSResolver dns.Resolver
var defaultLocationResolver location.Resolver

// rpcClient : grpc client.可直接使用,将在包初始化过程中完成初始化
var rpcClient GRPCClient

// GRPCClient :	grpc lient 对外接口
type GRPCClient interface {
	GetConnByAddr(remote net.Addr, local net.Addr) (net.Conn, error)
	Resolv(domain string) (ip string, err error)
}

// Init : 初始化 grpc client 池，解析器
func Init(serviceip string, serviceport int) (err error) {
	rpcClient, err = newstreamClient(serviceip, serviceport)
	if err != nil {
		logger.Error(ctxtransid.NewCTX(context.Background()), err, "server:", serviceip, serviceport)
	}
	defaultDNSResolver = dns.New(rpcClient.Resolv)
	defaultLocationResolver = location.New(resolvLocation)
	defaultPool, err = newPool(10, serviceip, int(serviceport), newclient)
	return
}

// Dial : 建立连接
func Dial(network string, ip string, port int) (net.Conn, error) {
	return dial(network, ip, port)
}

// ResolvDNS : 解析域名
func ResolvDNS(domain string) (ip string, err error) {
	return defaultDNSResolver.Resolv(domain)
}

// ResolvLocation : 解析 ip 所在地
func ResolvLocation(ip string) (loc locationbase.Location, err error) {
	loc, err = defaultLocationResolver.Resolv(ip)
	return
}
