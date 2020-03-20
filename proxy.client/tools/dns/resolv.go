/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-18 23:22:04
 * @LastEditTime: 2019-12-18 23:22:56
 */

package dns

import (
	"github.com/lancelot/proxy/proxy.client/tools/location"
	"github.com/lancelot/proxy/proxy.lib/dns"
	client "github.com/lancelot/proxy/proxy.server/grpc.client"
)

// Resolv 解析 DNS
var Resolv = ResolvByLocation

// ResolvByLocal 本地解析 DNS
func ResolvByLocal(domain string) (ip string, err error) {
	ip, err = dns.Resolv(domain)
	return
}

// ResolvByRemote 远程解析 DNS
func ResolvByRemote(domain string) (ip string, err error) {
	ip, err = client.ResolvDNS(domain)
	return
}

// ResolvByLocation 根据地址进行 DNS 解析
func ResolvByLocation(domain string) (ip string, err error) {
	ip, err = ResolvByLocal(domain)
	if err != nil {
		return
	}

	native, err := location.IsNativeIP(ip)
	if err != nil {
		return
	}

	if !native {
		ip, err = ResolvByRemote(domain)
	}

	return
}
