package limitconn

import (
	"context"
	"net"

	"github.com/pkg/errors"
	rate "golang.org/x/time/rate"
)

// ErrMaxSpeedisLessThanZero : 最大速度小于等于0
var ErrMaxSpeedisLessThanZero = errors.New("最大速度max小于等于零,应大于零")

// ErrBufferLenisLessThanZero :最大buffer长度小于等于0
var ErrBufferLenisLessThanZero = errors.New("缓存长度小于等于零,应大于零")

// limitConn : 带限速器的net.Conn
type limitConn struct {
	net.Conn  //底层net.Conn连接
	limiter   *rate.Limiter
	bufferLen int //每次传输的最大字节数，单位byte
}

// NewlimitConn : 构造NewlimitConn;
// base 表示底层的net.Conn;
// max表示每秒传输的最大字节数,单位Kbyte/s;
// bufferLen表示每次传输的最大字节数,单位Kbyte
func NewlimitConn(base net.Conn, max int, bufferLen int) (conn net.Conn, err error) {
	if max <= 0 {
		err = errors.WithStack(ErrMaxSpeedisLessThanZero)
		return
	}
	if bufferLen < 0 {
		err = errors.WithStack(ErrBufferLenisLessThanZero)
		return
	}
	maxCount := float64(max) / float64(bufferLen)
	conn = &limitConn{
		Conn:      base,
		bufferLen: bufferLen * 1000,
		limiter:   rate.NewLimiter(rate.Limit(maxCount), 3),
	}
	return
}

// Write : 实现net.Conn 的Write方法, 最多允许写入 bufferLen长度的数据
func (l *limitConn) Write(bf []byte) (n int, err error) {
	err = l.limiter.Wait(context.Background())
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(bf) > l.bufferLen {
		bf = bf[:l.bufferLen]
	}
	n, err = l.Conn.Write(bf)
	return
}

// Read : 实现 net.Conn 的Read方法, 最多允许读出 bufferLen长度的数据
func (l *limitConn) Read(bf []byte) (n int, err error) {
	err = l.limiter.Wait(context.Background())
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if len(bf) > l.bufferLen {
		bf = bf[:l.bufferLen]
	}
	n, err = l.Conn.Read(bf)
	return
}
