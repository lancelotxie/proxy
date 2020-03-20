package base

import (
	"crypto/sha1"
	"net"

	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
)

func dial(addr string) (conn net.Conn, err error) {
	key := "test pass"
	salt := "test salt"
	pass := pbkdf2.Key([]byte(key), []byte(salt), 4096, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(pass)

	s, err := kcp.DialWithOptions(addr, block, 10, 3)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer func() {
		if err != nil {
			s.Close()
		}
	}()

	s.SetStreamMode(true)
	s.SetWriteDelay(false)
	s.SetNoDelay(1, 10, 2, 1)
	s.SetWindowSize(1024, 1024)
	s.SetMtu(1350)
	s.SetACKNoDelay(true)

	if err = s.SetDSCP(46); err != nil {
		return
	}
	if err = s.SetReadBuffer(16777217); err != nil {
		return
	}
	if err = s.SetWriteBuffer(16777217); err != nil {
		return
	}

	conn = s
	return
}
