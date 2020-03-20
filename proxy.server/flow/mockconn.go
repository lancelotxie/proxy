package flow

import (
	"net"
	"sync"
	"time"

	errors "github.com/pkg/errors"
)

// BufferLen :
var BufferLen = 1000

// MockConn2Net : mock实现Conn2Net接口
type MockConn2Net struct {
	rch chan []byte
	wch chan []byte
}

// NewMockConn2Net :
func NewMockConn2Net() *MockConn2Net {
	return &MockConn2Net{
		rch: make(chan []byte, 1024),
		wch: make(chan []byte, 1024),
	}
}

// Trans :
func (mk *MockConn2Net) Trans(lconn net.Conn, rconn net.Conn) {
	waitg := &sync.WaitGroup{}
	waitg.Add(2)
	go func() {
		for {
			mk.left2right(lconn, rconn)
			time.Sleep(time.Millisecond * 100)
			if lconn == nil || rconn == nil {
				waitg.Done()
			}
		}
	}()
	go func() {
		for {
			mk.left2right(rconn, lconn)
			time.Sleep(time.Millisecond * 100)
			if lconn == nil || rconn == nil {
				waitg.Done()
			}
		}
	}()
	waitg.Wait()
}
func (mk *MockConn2Net) left2right(lconn net.Conn, rconn net.Conn) {
	rbytes := make([]byte, BufferLen)
	if lconn != nil {
		n, err := lconn.Read(rbytes)
		if err == nil && rconn != nil {
			rconn.Write(rbytes[:n])
		}
	}
}
func (mk *MockConn2Net) right2left(lconn net.Conn, rconn net.Conn) {
	rbytes := make([]byte, BufferLen)
	n, err := rconn.Read(rbytes)
	if err == nil {
		lconn.Write(rbytes[:n])
	}
}

// Close :
func (mk *MockConn2Net) Close() error { return nil }

// mockconn :模拟 一对相互读写的连接net.Conn
type mockconn struct {
}

// GetConn :模拟一对 互相读写的连接。lconn写入的数据，可以从 rconn读出数据；rconn写入数据，从lconn读出
func (m *mockconn) GetConn() (lconn net.Conn, rconn net.Conn) {
	leftch := make(chan []byte, 1024)
	rightch := make(chan []byte, 1024)
	lconn = &cusconn{
		rch: leftch, wch: rightch,
	}
	rconn = &cusconn{
		rch: rightch, wch: leftch,
	}
	return
}

// cusconn :实现net.Conn接口。从rch读数据。写数据到wch
type cusconn struct {
	closed bool
	wch    chan []byte
	rch    chan []byte
}

func (cc *cusconn) init() {
	cc.closed = false
	cc.wch = make(chan []byte, 1024)
	cc.rch = make(chan []byte, 1024)
}

// Write :
func (cc *cusconn) Write(bf []byte) (int, error) {
	if cc.closed {
		return 0, errors.New("conn was closed")
	}
	cc.wch <- bf
	return len(bf), nil
}

// Read :
func (cc *cusconn) Read(bf []byte) (int, error) {
	if cc.closed {
		return 0, errors.New("conn was closed")
	}
	data, ok := <-cc.rch
	if !ok {
		return 0, errors.New("chmsg no msg")
	}
	if cap(bf) < len(data) {
		return 0, errors.New("buffer is to short")
	}
	copy(bf, data)
	return len(data), nil
}

// Close :
func (cc *cusconn) Close() error {
	cc.closed = true
	return nil
}

// LocalAddr :
func (cc *cusconn) LocalAddr() net.Addr {
	return nil
}

// RemoteAddr :
func (cc *cusconn) RemoteAddr() net.Addr {
	return nil
}

// SetDeadline :
func (cc *cusconn) SetDeadline(to time.Time) error {
	return nil
}

// SetReadDeadline :
func (cc *cusconn) SetReadDeadline(to time.Time) error {
	return nil
}

// SetWriteDeadline :
func (cc *cusconn) SetWriteDeadline(to time.Time) error {
	return nil
}
