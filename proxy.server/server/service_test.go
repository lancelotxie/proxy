package service

import (
	"context"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	dial "github.com/lancelot/proxy/proxy.server/dial"
	flow "github.com/lancelot/proxy/proxy.server/flow"
	"github.com/lancelot/proxy/proxy.lib/address"
	"github.com/lancelot/proxy/proxy.lib/ctxtransid"
	"github.com/stretchr/testify/assert"
)

var BufferLen = 1024

func TestService(t *testing.T) {
	var (
		tunnel    flow.Conn2Net
		d2net     dial.Dial
		err       error
		connCount = 1000
	)
	// Init:
	mockflowA := &mockFlow{}
	mockflowA.Init()
	tunnel = mockflowA
	mockd2net := newmockDial2Net(connCount)
	d2net = mockd2net
	grpcListener := newmockgrpcListen(connCount)

	service := &server{
		base:   grpcListener,
		Dnet:   d2net,
		Tunnel: tunnel,
	}
	// Start:
	go service.serve(grpcListener)
	assert.Nil(t, err)
	// enter a conn
	addrA := address.New("144.144.144.144", 443, "tcp")
	readA := make(chan []byte, BufferLen)
	writeA := make(chan []byte, BufferLen)
	mockA := newmockconn(addrA, readA, writeA)
	assert.Zero(t, mockd2net.curidx)
	grpcListener.connOut <- mockA

	time.Sleep(time.Second) //等待服务端处理grpc过来的连接
	assert.NotZero(t, mockd2net.curidx)

}

func TestWrite(t *testing.T) {
	var (
		tunnel    flow.Conn2Net
		d2net     dial.Dial
		err       error
		connCount = 1000
	)
	tunnel = &flow.Conn2Conn{}
	mockd2net := newmockDial2Net(connCount)
	d2net = mockd2net
	grpcListener := newmockgrpcListen(connCount)

	service := &server{
		base:   grpcListener,
		Dnet:   d2net,
		Tunnel: tunnel,
	}
	// Start:
	go service.serve(grpcListener)
	assert.Nil(t, err)
	// enter a conn
	addrA := address.New("144.144.144.144", 443, "tcp")
	readA := make(chan []byte, BufferLen)
	writeA := make(chan []byte, BufferLen)
	mockA := newmockconn(addrA, readA, writeA)
	assert.Zero(t, mockd2net.curidx)
	grpcListener.connOut <- mockA

	time.Sleep(time.Second) //等待服务端处理grpc过来的连接
	assert.NotZero(t, mockd2net.curidx)
	mockB := mockd2net.lstconn[mockd2net.curidx-1]

	msgA := "client send"
	msgB := "web send"

	readA <- []byte(msgA)
	mockB.readfromchan <- []byte(msgB)
	outA := string(<-writeA)
	outB := string(<-mockB.writedata2chan)

	assert.Equal(t, msgA, outB)
	assert.Equal(t, msgB, outA)
}

type mockFlow struct {
	incount  int
	outcount int
	PQ       int64
	l        *sync.Mutex
	lstconn  []net.Conn
}

func (m *mockFlow) Init() {
	m.incount = 0
	m.outcount = 0
	m.PQ = 0
	m.lstconn = make([]net.Conn, 1000)
	m.l = &sync.Mutex{}
}

// 计数：从右边得到多少个net.conn
func (m *mockFlow) Trans(ctx context.Context, lconn net.Conn, rconn net.Conn) {
	log.Println("info:开始传输", m.incount)
	m.l.Lock()
	m.lstconn[m.incount] = rconn
	m.incount++
	m.l.Unlock()
	atomic.AddInt64(&m.PQ, 1)
	m.outcount++
}

func (m *mockFlow) Close() error { return nil }

type mockDial2Net struct {
	lstconn []*mockconn
	curidx  int32
}

func newmockDial2Net(size int) *mockDial2Net {
	return &mockDial2Net{
		lstconn: make([]*mockconn, size),
		curidx:  0,
	}
}

func (m *mockDial2Net) Conn2Web(ctx context.Context, remote net.Addr) (conn net.Conn, err error) {
	wch := make(chan []byte, BufferLen)
	rch := make(chan []byte, BufferLen)
	newmock := newmockconn(remote, wch, rch)
	num := atomic.LoadInt32(&m.curidx)
	atomic.AddInt32(&m.curidx, 1)
	m.lstconn[num] = newmock
	conn = newmock
	return
}

type mockconn struct {
	readfromchan   chan []byte
	writedata2chan chan []byte
	remote         net.Addr
	net.Conn
}

func newmockconn(remoteaddr net.Addr, rch chan []byte, wch chan []byte) *mockconn {
	return &mockconn{
		remote:         remoteaddr,
		readfromchan:   rch,
		writedata2chan: wch,
	}
}

func (m *mockconn) Read(bf []byte) (n int, err error) {
	data := <-m.readfromchan
	n = copy(bf, data)
	return
}

func (m *mockconn) Write(bf []byte) (n int, err error) {
	m.writedata2chan <- bf
	n = len(bf)
	return
}

func (m *mockconn) RemoteAddr() net.Addr { return m.remote }

type mockgrpcListen struct {
	curConn int32
	lstConn []*mockconn
	connOut chan *mockconn
}

func newmockgrpcListen(size int) *mockgrpcListen {
	return &mockgrpcListen{
		curConn: 0,
		lstConn: make([]*mockconn, size),
		connOut: make(chan *mockconn, size),
	}
}

func (m *mockgrpcListen) Accept() (ctx context.Context, conn net.Conn, err error) {
	ctx = ctxtransid.NewCTX(context.Background())
	conn = <-m.connOut
	return
}
func (m *mockgrpcListen) Close() error   { return nil }
func (m *mockgrpcListen) Addr() net.Addr { return nil }
