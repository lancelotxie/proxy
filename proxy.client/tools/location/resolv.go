package location

import (
	"github.com/lancelotXie/proxy/proxy.lib/location/base"
	client "github.com/lancelotXie/proxy/proxy.server/grpc.client"

	mynet "github.com/eaglexiang/go/net"
)

// 本程序所在地
var localLocation = base.China

// resolv 解析 IP 的 地理位置
func resolv(ip string) (loc base.Location, err error) {
	loc, err = client.ResolvLocation(ip)
	return
}

// isNative 判断 目标所在地 是否为 本地
func isNative(location base.Location) (ok bool) {
	if location == localLocation {
		ok = true
	}
	return
}

// isLANIP 判断 IP 是否属于 LAN
func isLANIP(ip string) (ok bool) {
	ok = mynet.IsPrivateIPv4(ip)
	return
}

// IsNativeIP 判断是否是本土 IP
func IsNativeIP(ip string) (ok bool, err error) {
	if ok = isLANIP(ip); ok {
		return
	}

	loc, err := resolv(ip)
	if err != nil {
		return
	}

	ok = isNative(loc)
	return
}
