package flow

import (
	"context"
	"net"

	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
	tunnel "github.com/eaglexiang/go/tunnel"
)

// Conn2Net : 连接2个net.Conn
type Conn2Net interface {
	Close() (err error)
	Trans(context.Context, net.Conn, net.Conn)
}

// Conn2Conn : 连通2个net.Conn
type Conn2Conn struct {
	leftconn  net.Conn
	rightconn net.Conn
}

// Close :
func (c *Conn2Conn) Close() error { return nil }

// Trans : 连通2个net.Conn
func (c *Conn2Conn) Trans(ctx context.Context, lconn net.Conn, rconn net.Conn) {
	if lconn == nil {
		logger.Error(ctx, "左边连接为nil")
		return
	}
	if rconn == nil {
		logger.Error(ctx, "右边连接为nil")
	}
	c.leftconn = lconn
	c.rightconn = rconn
	t := tunnel.GetTunnel()
	t.SetLeft(lconn)
	t.SetRight(rconn)
	logger.Info(ctx, "Tunnel 开始流通")
	id := ctxtransid.GetID(ctx)
	t.Flow()
	t.Close()
	tunnel.PutTunnel(t)
	logger.Info(context.Background(), "id:", id, "Tunnel 停止流通")
}
