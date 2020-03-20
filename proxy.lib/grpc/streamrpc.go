package grpc

import (
	"errors"
	"net"
	"time"
)

// RPCConn :
type RPCConn struct {
}

// New :
func New() *RPCConn {
	return &RPCConn{}
}

// Close :
func (cc *RPCConn) Close() error {
	return nil
}

// LocalAddr :
func (cc *RPCConn) LocalAddr() net.Addr {
	return nil
}

// Read :
func (cc *RPCConn) Read(bf []byte) (int, error) {
	return 0, errors.New("chmsg no msg")
}

// RemoteAddr :
func (cc *RPCConn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline :
func (cc *RPCConn) SetDeadline(to time.Time) error {
	return nil
}

// SetReadDeadline :
func (cc *RPCConn) SetReadDeadline(to time.Time) error {
	return nil
}

// SetWriteDeadline :
func (cc *RPCConn) SetWriteDeadline(to time.Time) error {
	return nil
}

// Write :
func (cc *RPCConn) Write(bf []byte) (int, error) {

	return len(bf), nil
}
