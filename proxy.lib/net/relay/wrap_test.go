package relay

import (
	"bytes"
	"context"
	"net"
	"testing"

	mocknet "github.com/lancelot/proxy/proxy.lib/mock/net"

	"github.com/stretchr/testify/assert"
)

func Test_tcpRelay(t *testing.T) {
	ctx := context.Background()
	dstAddr := "127.0.0.1:8899"

	dstLis, err := net.Listen("tcp", dstAddr)
	if !assert.NoError(t, err) {
		return
	}

	done := make(chan struct{})
	go func() {
		conn, err := dstLis.Accept()
		if !assert.NoError(t, err) {
			return
		}

		buf := bytes.NewBuffer(nil)
		_, err = buf.ReadFrom(conn)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "test msg", buf.String())
		close(done)
	}()

	src, dial := mocknet.NewListener()
	s := NewTCP(ctx, dstAddr, src)
	go s.Start()

	conn, err := dial()
	if !assert.NoError(t, err) {
		return
	}

	buf := bytes.NewBufferString("test msg")
	_, err = buf.WriteTo(conn)
	assert.NoError(t, err)

	conn.Close()

	<-done
}
