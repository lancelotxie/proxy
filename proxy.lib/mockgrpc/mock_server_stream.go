package mockgrpc

import (
	"context"

	proto "github.com/lancelot/proxy/proxy.server/grpc.server/proto"

	"google.golang.org/grpc/metadata"
)

type serverStreamInterface interface {
	proto.StreamService_StreamDualServer
}

// serverstream 服务端使用的 流
type serverstream struct {
	Clientstream  proto.StreamService_StreamDualClient //暂时没用
	client        proto.StreamServiceClient            //暂时没用
	server2client chan proto.StreamBytes               //单向流，服务端到客户端
	client2server chan proto.StreamBytes               //单向流，客户端到服务端
}

func newserverstream(
	gclient proto.StreamServiceClient,
	clientStream proto.StreamService_StreamDualClient,
	server2c chan proto.StreamBytes,
	client2s chan proto.StreamBytes) *serverstream {
	return &serverstream{
		Clientstream:  clientStream,
		client:        gclient,
		server2client: server2c,
		client2server: client2s,
	}
}

// 发送数据 到 channel, 客户端使用这个channel接收数据.前提是这个channel 需要提前给client设置好
func (d *serverstream) Send(data *proto.StreamBytes) error {
	d.server2client <- *data
	return nil
}

// 从channel接收数据, 客户端将写入这个channel。前提是这个channel 需要提前给client设置好
func (d *serverstream) Recv() (res *proto.StreamBytes, err error) {
	data := <-d.client2server
	res = &data
	return
}

func (d *serverstream) SetHeader(metadata.MD) error  { return nil }
func (d *serverstream) SendHeader(metadata.MD) error { return nil }
func (d *serverstream) SetTrailer(metadata.MD)       {}
func (d *serverstream) Context() context.Context     { return nil }
func (d *serverstream) SendMsg(m interface{}) error  { return nil }
func (d *serverstream) RecvMsg(m interface{}) error  { return nil }
