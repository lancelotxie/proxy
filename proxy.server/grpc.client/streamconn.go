package client

import (
	"bytes"
	"context"
	"net"
	"time"

	"github.com/lancelotXie/proxy/proxy.lib/logger"
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"

	errors "github.com/pkg/errors"
)

var bufferLen = 64_000

type clientConn struct {
	ctx        context.Context
	localaddr  net.Addr
	remoteaddr net.Addr
	cancel     context.CancelFunc
	stream     proto.StreamService_StreamDualClient
	buf        *bytes.Buffer
}

func newclientConn(ctx context.Context, local net.Addr, remote net.Addr, c context.CancelFunc,
	s proto.StreamService_StreamDualClient) *clientConn {
	return &clientConn{
		ctx:        ctx,
		localaddr:  local,
		remoteaddr: remote,
		cancel:     c,
		stream:     s,
		buf:        new(bytes.Buffer),
	}
}

// Read :
func (c *clientConn) Read(b []byte) (n int, err error) {
	if c.buf.Len() > 0 {
		n, err = c.buf.Read(b)
		return
	}

	rep, err := c.stream.Recv()
	if err != nil {
		return
	}

	c.buf = bytes.NewBuffer(rep.Data)
	n, err = c.buf.Read(b)
	return
}

func (c *clientConn) Write(b []byte) (n int, err error) {
	buf := bytes.NewBuffer(b)
	for data := buf.Next(bufferLen); len(data) > 0; data = buf.Next(bufferLen) {
		sendreq := &proto.StreamBytes{
			Data: b,
		}
		err = c.stream.Send(sendreq)
		n += len(data)
	}
	if n != len(b) {
		err = errors.New("err:write to long")
		logger.Error(c.ctx, err)
	}
	return
}

// Close :
func (c *clientConn) Close() error {
	logger.Info(c.ctx, "关闭连接")
	c.cancel()
	return nil
}

// LocalAddr :
func (c *clientConn) LocalAddr() net.Addr { return c.localaddr }

// RemoteAddr :
func (c *clientConn) RemoteAddr() net.Addr { return c.remoteaddr }

// SetDeadline :
func (c *clientConn) SetDeadline(t time.Time) error { return nil }

// SetReadDeadline :
func (c *clientConn) SetReadDeadline(t time.Time) error { return nil }

// SetWriteDeadline :
func (c *clientConn) SetWriteDeadline(t time.Time) error { return nil }
