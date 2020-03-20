/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-05 22:02:15
 * @LastEditTime: 2019-12-07 16:16:38
 */

package https

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

	"github.com/pkg/errors"
)

func init() {
	h := &HTTPS{}
	parsers.RegisterParser(h)
}

// re443 CONNECT 方法成功的返回值
const re443 = "HTTP/1.1 200 Connection Established\r\n\r\n"

// ErrHeadNotMatched 头部标记未能找到
var ErrHeadNotMatched = errors.New("head not matched")

// ErrFailed2ReplyClient 回复 HTTP 客户端失败
var ErrFailed2ReplyClient = errors.New("failed to reply http client")

// HTTPS HTTPS 流量解析器
type HTTPS struct {
}

// Match 判断流量是否匹配
func (h HTTPS) Match(b []byte) (err error) {
	head := string(b)
	items := strings.Split(head, " ")
	if len(items) < 2 {
		err = errors.WithStack(errors.WithMessage(ErrHeadNotMatched, head))
		return
	}
	method := items[0]

	if method != "CONNECT" {
		err = errors.WithStack(errors.WithMessage(ErrHeadNotMatched, method))
	}
	return
}

// Parse 解析流量
func (h HTTPS) Parse(ctx context.Context, connOfReq net.Conn) (conn2Remote net.Conn, err error) {
	reader := bufio.NewReader(connOfReq)
	req, err := http.ReadRequest(reader)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	host := req.URL.Hostname()

	port := req.URL.Port()
	_port, err := fixPort(port)
	if err != nil {
		return
	}

	log.Println("捕获请求：", host, ":", _port)
	conn2Remote, err = dial.TCPDial(ctx, host, _port)
	if err != nil {
		return
	}

	bytesRe := []byte(re443)
	n, err := connOfReq.Write(bytesRe)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if n != len(bytesRe) {
		err = errors.WithStack(ErrFailed2ReplyClient)
	}
	return
}

func (h HTTPS) String() string {
	return "HTTPS"
}

func fixPort(port string) (dst int, err error) {
	if port == "" {
		dst = 443
		return
	}

	dst, err = strconv.Atoi(port)
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}
