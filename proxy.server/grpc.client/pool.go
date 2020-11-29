package client

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
	errors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

// Errclosed :
var Errclosed = errors.New("already closed")

// ErrOverflow :
var ErrOverflow = errors.New("Overflow")

// ErrNil :
var ErrNil = errors.New("nil")

type client struct {
	conn *grpc.ClientConn
	t    time.Timer
}

type pool struct {
	capacity int
	lstconn  []proto.StreamServiceClient
	locker   *sync.Mutex
	putnext  int32
	factory  func(string, int) (*grpc.ClientConn, error)
	ip       string
	port     int
}

func newPool(size int, ip string, port int, factory func(string, int) (*grpc.ClientConn, error)) (newp *pool, err error) {
	newp = &pool{
		capacity: size,
		putnext:  0,
		locker:   &sync.Mutex{},
		lstconn:  make([]proto.StreamServiceClient, size),
		factory:  factory,
		ip:       ip,
		port:     port,
	}
	for idx := range newp.lstconn {
		newconn, err2 := factory(ip, port)
		if err2 != nil {
			err = err2
			return
		}
		client := proto.NewStreamServiceClient(newconn)
		newp.lstconn[idx] = client
	}
	return
}

func (p *pool) GetClient() (conn proto.StreamServiceClient, err error) {
	index := atomic.LoadInt32(&p.putnext) % int32(p.capacity)
	atomic.AddInt32(&p.putnext, 1)
	log.Println("坐标：", index)
	conn = p.lstconn[index]
	return
}
