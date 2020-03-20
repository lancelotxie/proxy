package mockgrpc

import (
	"context"
	"sync"

	proto "github.com/lancelot/proxy/proxy.server/grpc.server/proto"

	"google.golang.org/grpc/metadata"
)

type clientStreamInterface interface {
	proto.StreamService_StreamDualClient
}

// clientstream : 客户端使用的 流
type clientstream struct {
	enterDualService *sync.Once                           //标示 仅进入一次 server端的DualStream接口方法
	nohupDualService bool                                 //标示 是否后台运行 处理 server端的DualStream接口方法
	Serverstream     proto.StreamService_StreamDualServer //可以互传消息的 服务端的流
	server           proto.StreamServiceServer            //连接的服务端
	server2client    chan proto.StreamBytes               //单向流，服务端到客户端
	client2server    chan proto.StreamBytes               //单向流，客户端到服务端
	sendDo           func(*clientstream)                  //第一次发送时，运行的方法
}

func newclientstream(
	gserver proto.StreamServiceServer,
	serverStream proto.StreamService_StreamDualServer,
	client2s chan proto.StreamBytes,
	server2c chan proto.StreamBytes,
	nohupServer bool) *clientstream {
	return &clientstream{
		enterDualService: &sync.Once{},
		Serverstream:     serverStream,
		server:           gserver,
		server2client:    server2c,
		client2server:    client2s,
		nohupDualService: nohupServer,
		sendDo:           serverDual,
	}
}

func serverDual(d *clientstream) {
	if d.nohupDualService {
		go d.server.StreamDual(d.Serverstream)
	} else {
		d.server.StreamDual(d.Serverstream)
	}
}

// 发送数据 到 channel, 服务端使用这个channel接收数据.前提是这个channel 需要提前给server设置好
func (d *clientstream) Send(data *proto.StreamBytes) error {
	d.client2server <- *data
	d.enterDualService.Do(func() {
		if d.sendDo != nil {
			d.sendDo(d)
		}
	})
	return nil
}

// 从channel接收数据, 服务端将写入这个channel。前提是这个channel 需要提前给server设置好
func (d *clientstream) Recv() (res *proto.StreamBytes, err error) {
	data := <-d.server2client
	res = &data
	return
}

func (d *clientstream) Header() (metadata.MD, error) { return nil, nil }
func (d *clientstream) Trailer() metadata.MD         { return nil }
func (d *clientstream) CloseSend() error             { return nil }
func (d *clientstream) Context() context.Context     { return nil }
func (d *clientstream) SendMsg(m interface{}) error  { return nil }
func (d *clientstream) RecvMsg(m interface{}) error  { return nil }
