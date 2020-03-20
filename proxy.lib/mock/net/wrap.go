package net

import "net"

// NewListener 构造新的 mock Listenr
func NewListener() (lis net.Listener, dialer DialFunc) {
	l, dialer := newListener()
	lis = l
	return
}
