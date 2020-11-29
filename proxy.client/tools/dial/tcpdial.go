/*
 * @Author: EagleXiang
 * @LastEditors  : EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 16:04:35
 * @LastEditTime : 2019-12-18 22:26:02
 */

package dial

import (
	"context"
	"fmt"
	"net"

	"github.com/lancelotXie/proxy/proxy.client/tools/dns"
	"github.com/lancelotXie/proxy/proxy.client/tools/location"
	client "github.com/lancelotXie/proxy/proxy.server/grpc.client"

	"github.com/pkg/errors"
)

// TCPDial 建立 TCP 连接
var TCPDial = func(ctx context.Context, ip string, port int) (conn net.Conn, err error) {
	return tcpDialByLocation(ctx, ip, port)
}

// tcpDialByLocal 通过本地建立 TCP 连接
func tcpDialByLocal(ip string, port int) (conn net.Conn, err error) {
	address := fmt.Sprint(ip, ":", port)
	conn, err = net.Dial("tcp", address)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

// 利用远端 建立连接
func tcpDialByRemote(ip string, port int) (conn net.Conn, err error) {
	conn, err = client.Dial("tcp", ip, port)
	return
}

// 根据ip所在地，分别用 本地建立连接 或者 远端连接
func tcpDialByLocation(ctx context.Context, ip string, port int) (conn net.Conn, err error) {
	ip, err = dns.ResolvByLocation(ip)
	if err != nil {
		return
	}
	native, err := location.IsNativeIP(ip)
	if err != nil {
		return
	}

	if native {
		conn, err = tcpDialByLocal(ip, port)
	} else {
		conn, err = tcpDialByRemote(ip, port)
	}

	return
}
