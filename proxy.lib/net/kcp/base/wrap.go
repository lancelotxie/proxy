package base

import (
	"net"
)

// Dial 拨号
func Dial(addr string) (conn net.Conn, err error) {
	conn, err = dial(addr)
	return
}

// Listen 监听
func Listen(addr string) (lis net.Listener, err error) {
	l, err := listen(addr)
	if err != nil {
		return
	}

	lis = l
	return
}
