package relay

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"github.com/eaglexiang/go/tunnel"
	"github.com/pkg/errors"
)

// DialFunc 拨号函数
type DialFunc func() (conn net.Conn, err error)

type relay struct {
	ctx       context.Context
	dial      DialFunc
	src       net.Listener
	conns2DST <-chan net.Conn
	cancel    context.CancelFunc
}

func newRelay(ctx context.Context, dialer DialFunc, src net.Listener) (r *relay) {
	r = new(relay)
	r.ctx = ctx
	r.dial = dialer
	r.src = src
	return
}

// 开始服务
func (r *relay) Start() (err error) {
	conns2DST := make(chan net.Conn)
	r.conns2DST = conns2DST

	ctx, cancel := context.WithCancel(r.ctx)
	r.cancel = cancel

	go r.produceConns(ctx, conns2DST)
	go r.listen(ctx)

	return
}

// listen 监听转发请求
func (r *relay) listen(ctx context.Context) {
	log.Println("开始监听转发请求")
	for {
		conn, err := r.src.Accept()
		if err != nil {
			break
		}

		go r.handle(ctx, conn)
	}
	log.Println("停止监听转发请求")
}

// produceConnss 预生成连接
func (r *relay) produceConns(ctx context.Context, conns2DST chan<- net.Conn) {
	log.Println("开始预建立转发连接")
loop:
	for {
		conn, err := r.dial()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}

		select {
		case conns2DST <- conn:
		case <-ctx.Done():
			close(conns2DST)
			break loop
		}
	}
	log.Println("停止预建立转发连接")
}

// handle 获取面向目的地址的连接，并建立与 src 之间的双向流动
func (r *relay) handle(ctx context.Context, src net.Conn) {
	var dst net.Conn
	var ok bool

	select {
	case dst, ok = <-r.conns2DST:
	case <-ctx.Done():
		return
	}
	if !ok {
		err := errors.WithStack(io.EOF)
		log.Println(err)
		return
	}

	tunnel.Flow(src, dst)
}

// Close 关闭服务
func (r *relay) Close() (err error) {
	err = r.src.Close()
	r.cancel()
	return
}
