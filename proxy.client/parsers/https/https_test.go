/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-05 20:38:57
 * @LastEditTime: 2019-12-07 15:42:44
 */

package https

import (
	"context"
	"net"
	"testing"

	"github.com/lancelot/proxy/proxy.client/tools/dial"
	"github.com/lancelot/proxy/proxy.client/tools/dns"

	"github.com/eaglexiang/go/tunnel"
	"github.com/stretchr/testify/assert"
)

func init() {
	// mock dns resolver
	dns.Resolv = func(domain string) (ip string, err error) {
		ip = domain
		return
	}
	// mock tcp dialer
	dial.TCPDial = func(ctx context.Context, ip string, port int) (conn net.Conn, err error) {
		conn = tunnel.NewVirtualConn()
		return
	}
}

func Test_HTTPS_Match(t *testing.T) {
	h := HTTPS{}

	b := []byte("CONNECT")
	err := h.Match(b)
	assert.Error(t, err)

	b = []byte("CONNECT AAAA")
	err = h.Match(b)
	assert.NoError(t, err)

	b = []byte("GET HHHH")
	err = h.Match(b)
	assert.Error(t, err)
}

func Test_Parse(t *testing.T) {
	var req = "CONNECT www.bilibili.com:443 HTTP/1.1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:70.0) Gecko/20100101 Firefox/70.0\r\nProxy-Connection: keep-alive\r\nConnection: keep-alive\r\nHost: www.bilibili.com:443\r\n\r\n"

	mockConn := tunnel.NewVirtualConn()
	mockConn.Write([]byte(req))

	ctx := context.Background()
	h := HTTPS{}
	_, err := h.Parse(ctx, mockConn)
	assert.NoError(t, err)
}

func Test_fixPort(t *testing.T) {
	port := "80"
	p, err := fixPort(port)
	assert.NoError(t, err)
	assert.Equal(t, 80, p)

	port = ""
	p, err = fixPort(port)
	assert.NoError(t, err)
	assert.Equal(t, 443, p)

	port = "a"
	_, err = fixPort(port)
	assert.Error(t, err)
}
