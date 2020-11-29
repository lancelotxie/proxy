package relay

import (
	"context"
	"net"

	mynet "github.com/lancelotXie/proxy/proxy.lib/net"
)

// New 构造新的流量中继器
func New(ctx context.Context, dialer DialFunc, src net.Listener) (s mynet.Service) {
	r := newRelay(ctx, dialer, src)
	s = r
	return
}

// NewTCP 构造新的TCP流量中继器
func NewTCP(ctx context.Context, dstAddr string, src net.Listener) (s mynet.Service) {
	dialer := func() (conn net.Conn, err error) {
		conn, err = net.Dial("tcp", dstAddr)
		return
	}

	r := newRelay(ctx, dialer, src)
	s = r
	return
}
