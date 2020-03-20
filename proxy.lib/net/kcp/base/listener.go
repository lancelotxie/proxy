package base

import (
	"crypto/sha1"
	"net"

	"github.com/pkg/errors"
	"golang.org/x/crypto/pbkdf2"

	"github.com/xtaci/kcp-go"
)

type listener struct {
	base *kcp.Listener
}

func listen(addr string) (lis *listener, err error) {
	key := "test pass"
	salt := "test salt"
	pass := pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(pass)

	base, err := kcp.ListenWithOptions(addr, block, 10, 3)
	if err != nil {
		err = errors.WithStack(err)
	}
	defer func() {
		if err != nil {
			base.Close()
		}
	}()

	if err = base.SetDSCP(46); err != nil {
		return
	}
	if err = base.SetReadBuffer(16777217); err != nil {
		return
	}
	if err = base.SetWriteBuffer(16777217); err != nil {
		return
	}

	lis = new(listener)
	lis.base = base
	return
}

// Accept 接收新的连接
func (l *listener) Accept() (conn net.Conn, err error) {
	s, err := l.base.AcceptKCP()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	s.SetStreamMode(true)
	s.SetWriteDelay(false)
	s.SetNoDelay(1, 10, 2, 1)
	s.SetMtu(1350)
	s.SetWindowSize(1024, 1024)
	s.SetACKNoDelay(true)

	conn = s
	return
}

func (l *listener) Close() (err error) {
	err = l.base.Close()
	return
}

func (l *listener) Addr() (addr net.Addr) {
	addr = l.base.Addr()
	return
}
