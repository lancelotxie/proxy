/*
 * @Author: EagleXiang
 * @LastEditors  : EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 13:24:22
 * @LastEditTime : 2019-12-18 22:24:48
 */

package http

import (
	"bufio"
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/lancelot/proxy/proxy.client/parsers"
	"github.com/lancelot/proxy/proxy.client/tools/dial"
	"github.com/lancelot/proxy/proxy.lib/container/set"

	"github.com/pkg/errors"
)

// ErrHeadNotMatched 头部标记未能找到
var ErrHeadNotMatched = errors.New("head not matched")

// ErrInvalidScheme 非法协议
var ErrInvalidScheme = errors.New("invalid scheme")

var validMethods *set.Set

func init() {
	validMethods = set.NewSet()
	validMethods.Set("OPTIONS")
	validMethods.Set("HEAD")
	validMethods.Set("GET")
	validMethods.Set("POST")
	validMethods.Set("PUT")
	validMethods.Set("DELETE")
	validMethods.Set("TRACE")
}

func init() {
	h := &HTTP{}
	parsers.RegisterParser(h)
}

// HTTP HTTP 解析器
type HTTP struct {
}

// Match 判断流量是否匹配
func (h HTTP) Match(b []byte) (err error) {
	head := string(b)
	items := strings.Split(head, " ")
	if len(items) < 2 {
		err = errors.WithStack(errors.WithMessage(ErrHeadNotMatched, head))
		return
	}
	method := items[0]

	valid := validMethods.Exist(method)
	if !valid {
		err = errors.WithStack(errors.WithMessage(ErrHeadNotMatched, method))
	}
	return
}

// Parse 解析流量
func (h HTTP) Parse(ctx context.Context, connOfReq net.Conn) (conn2Remote net.Conn, err error) {
	reader := bufio.NewReader(connOfReq)
	req, err := http.ReadRequest(reader)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	newReq, err := rebuildRequest(req)
	if err != nil {
		return
	}

	scheme := newReq.URL.Scheme

	host := newReq.URL.Hostname()
	host = fixHost(host)

	port := newReq.URL.Port()
	_port, err := fixPort(scheme, port)
	if err != nil {
		return
	}

	log.Println("捕获请求：", scheme, host, ":", _port)
	conn2Remote, err = dial.TCPDial(ctx, host, _port)
	if err != nil {
		return
	}

	err = newReq.Write(conn2Remote)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (h HTTP) String() string {
	return "HTTP"
}

// rebuildRequest 重建请求，以保证兼容性，以及隐匿代理服务的存在
func rebuildRequest(src *http.Request) (dst *http.Request, err error) {
	method := src.Method
	url := src.URL.String()
	body := src.Body

	dst, err = http.NewRequest(method, url, body)
	if err != nil {
		err = errors.WithStack(err)
	}

	copyHeader(src, dst)

	return
}

func fixHost(host string) string {
	items := strings.Split(host, ":")
	return items[0]
}

func fixPort(scheme, port string) (dst int, err error) {
	if port != "" {
		dst, err = strconv.Atoi(port)
		if err != nil {
			err = errors.WithStack(errors.WithMessage(err, port))
			return
		}
		return
	}

	scheme = strings.ToLower(scheme)
	switch scheme {
	case "http":
		dst = 80
	case "https":
		dst = 443
	default:
		err = errors.WithStack(ErrInvalidScheme)
	}
	return
}

func copyHeader(src, dst *http.Request) {
	headers := src.Header

	// closeKeepAlive(headers)
	removeProxyHeader(headers)

	for key, values := range headers {
		for _, value := range values {
			dst.Header.Set(key, value)
		}
	}
}

func closeKeepAlive(headers http.Header) {
	headers.Del("Connection")
	headers.Del("connection")
	headers.Set("Connection", "close")
}

func removeProxyHeader(headers http.Header) {
	headers.Del("Proxy-Connection")
	headers.Del("proxy-connection")
}
