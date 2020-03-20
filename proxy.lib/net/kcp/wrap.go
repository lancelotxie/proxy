package kcp

import "net"

// Dial 拨号
func Dial(addr string) (conn net.Conn, err error) {
	conn, err = dial(addr)
	return
}

// Listen 监听地址
func Listen(addr string) (lis net.Listener, err error) {
	l, err := listen(addr)
	lis = l
	return
}
