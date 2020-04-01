package server

import (
	"bytes"
	"context"
	"net"
	"time"

	proto "github.com/lancelot/proxy/proxy.server/grpc.server/proto"
	errors "github.com/pkg/errors"
)

var bufferLen = 64_000

// streamConn :
type streamConn struct {
	ctx        context.Context
	localaddr  net.Addr
	remoteaddr net.Addr
	stream     proto.StreamService_StreamDualServer
	c          context.CancelFunc
	buf        *bytes.Buffer
}

// GetAConn :
func GetAConn(ctx context.Context, cstream proto.StreamService_StreamDualServer) net.Conn {
	return newstreamConn(ctx, cstream, nil, nil, nil)
}

// NewStreamConn : construct NewStreamConn
func newstreamConn(ctx context.Context, stream proto.StreamService_StreamDualServer,
	cancel context.CancelFunc, remote net.Addr, local net.Addr) *streamConn {
	return &streamConn{
		ctx:        ctx,
		stream:     stream,
		c:          cancel,
		remoteaddr: remote,
		localaddr:  local,
		buf:        new(bytes.Buffer),
	}
}

// Read :
func (s *streamConn) Read(bf []byte) (n int, err error) {
	if s.buf.Len() > 0 {
		n, err = s.buf.Read(bf)
		return
	}

	recv, err := s.stream.Recv()
	if err != nil {
		return
	}

	s.buf = bytes.NewBuffer(recv.Data)
	n, err = s.buf.Read(bf)
	return
}

// Write :
func (s *streamConn) Write(bf []byte) (n int, err error) {
	return write(bf, s.stream)
}

// Close :
func (s *streamConn) Close() error {
	s.c()
	return nil
}

// LocalAddr :
func (s *streamConn) LocalAddr() net.Addr { return s.localaddr }

// RemoteAddr :
func (s *streamConn) RemoteAddr() net.Addr { return s.remoteaddr }

// SetDeadline :
func (s *streamConn) SetDeadline(t time.Time) error { return nil }

// SetReadDeadline :
func (s *streamConn) SetReadDeadline(t time.Time) error { return nil }

// SetWriteDeadline :
func (s *streamConn) SetWriteDeadline(t time.Time) error { return nil }

func read(stream proto.StreamService_StreamDualServer, bf []byte) (n int, err error) {
	recv, err := stream.Recv()
	if err != nil {
		return
	}
	if len(recv.Data) > cap(bf) {
		err = errors.New("err:data is too long")
		return
	}
	n = copy(bf, recv.Data)
	return
}

func write(bf []byte, stream proto.StreamService_StreamDualServer) (n int, err error) {
	buf := bytes.NewBuffer(bf)
	for smallbuf := buf.Next(bufferLen); len(smallbuf) > 0; smallbuf = buf.Next(bufferLen) {
		data := &proto.StreamBytes{
			Data: smallbuf,
		}
		err = stream.Send(data)
		if err != nil {
			return
		}
		n += len(smallbuf)
	}
	return
}
