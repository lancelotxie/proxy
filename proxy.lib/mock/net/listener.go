package net

import (
	"io"
	"net"

	"github.com/pkg/errors"
)

// DialFunc 拨号方法
type DialFunc func() (conn net.Conn, err error)

type listener struct {
	conns chan net.Conn
}

func newListener() (lis *listener, dialer DialFunc) {
	lis = new(listener)
	lis.conns = make(chan net.Conn, 10)

	dialer = func() (conn net.Conn, err error) {
		a, b := net.Pipe()
		conn = a
		lis.conns <- b
		return
	}

	return
}

func (lis *listener) Accept() (conn net.Conn, err error) {
	conn, ok := <-lis.conns
	if !ok {
		err = errors.WithStack(io.EOF)
	}
	return
}

func (lis *listener) Addr() (addr net.Addr) {
	return
}

func (lis *listener) Close() (err error) {
	return
}
