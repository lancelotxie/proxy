/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 13:40:08
 * @LastEditTime: 2019-12-07 15:32:17
 */

package http

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

func Test_Math(t *testing.T) {
	h := HTTP{}

	head := "GET"
	err := h.Match([]byte(head))
	assert.Error(t, err)

	head = "CONNECT AFA"
	err = h.Match([]byte(head))
	assert.Error(t, err)

	head = "OPTIONS ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)

	head = "HEAD ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)

	head = "GET ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)

	head = "POST ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)

	head = "PUT ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)

	head = "DELETE ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)

	head = "TRACE ALLFM"
	err = h.Match([]byte(head))
	assert.NoError(t, err)
}

func Test_Parse(t *testing.T) {
	var req = "GET http://10086.cn/ HTTP/1.1\r\nHost: 10086.cn\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:70.0) Gecko/20100101 Firefox/70.0\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8\r\nAccept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2\r\nAccept-Encoding: gzip, deflate\r\nConnection: keep-alive\r\nUpgrade-Insecure-Requests: 1\r\n\r\n"

	mockConn := tunnel.NewVirtualConn()
	mockConn.Write([]byte(req))

	ctx := context.Background()
	h := HTTP{}
	_, err := h.Parse(ctx, mockConn)
	assert.NoError(t, err)
}

func Test_fixPort(t *testing.T) {
	scheme := "http"
	port := "20"
	p, err := fixPort(scheme, port)
	assert.NoError(t, err)
	assert.Equal(t, 20, p)

	port = ""
	p, err = fixPort(scheme, port)
	assert.NoError(t, err)
	assert.Equal(t, 80, p)

	scheme = "https"
	p, err = fixPort(scheme, port)
	assert.NoError(t, err)
	assert.Equal(t, 443, p)

	port = "a"
	_, err = fixPort(scheme, port)
	assert.Error(t, err)

	scheme = "htt"
	port = ""
	_, err = fixPort(scheme, port)
	assert.Error(t, err)
}
