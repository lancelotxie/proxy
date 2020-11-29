/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-01 14:26:15
 * @LastEditTime: 2019-12-07 16:50:13
 */

package parsers

import (
	"context"
	"log"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/renewconn"

	"github.com/pkg/errors"
)

// ErrParserNotFound 没有找到匹配的解析器
var ErrParserNotFound = errors.New("parser not found")

// ErrParserExist 解析器已存在
var ErrParserExist = errors.New("parser exist")

// HeaderBufLen header-buf 的长度
const HeaderBufLen = 1024

// Parser 请求解析器
type Parser interface {
	Match([]byte) error
	Parse(ctx context.Context, conn net.Conn) (net.Conn, error)
	String() string
}

var allParsers map[string]Parser

func init() {
	allParsers = make(map[string]Parser)
}

func getHeader(conn net.Conn) (buf []byte, err error) {
	b := make([]byte, HeaderBufLen)
	count, err := conn.Read(b)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	buf = b[:count]
	return
}

func matchParser(header []byte) (p Parser, err error) {
	for _, tmpP := range allParsers {
		err = tmpP.Match(header)
		if err == nil {
			p = tmpP
			log.Println("命中解析器：", p.String())
			return
		}
	}
	err = errors.WithStack(ErrParserNotFound)
	return
}

// GetParser 匹配一个合适的解析器，返回值 renew 的右边包含完整的数据流拷贝
func GetParser(conn net.Conn) (p Parser, renew net.Conn, err error) {
	header, err := getHeader(conn)
	if err != nil {
		return
	}

	p, err = matchParser(header)
	if err != nil {
		return
	}

	renew = renewconn.Renew(conn, header)
	return
}

// RegisterParser 注册一个流量解析器
func RegisterParser(p Parser) (err error) {
	name := p.String()

	_, ok := allParsers[name]
	if ok {
		log.Println("解析器注册失败：", name, "已存在")
		err = errors.WithStack(ErrParserExist)
	} else {
		allParsers[p.String()] = p
		log.Println("解析器注册成功：", name)
	}
	return
}
