package mockgrpc

import (
	"context"
	"net"

	proto "github.com/lancelot/proxy/proxy.server/grpc.server/proto"
	"google.golang.org/grpc"
)

type grpcClientInterface interface {
	proto.StreamServiceClient
}

// mockClient : 模拟grpc 客户端
type mockClient struct {
	testgrpcServer bool                      //辨别测试 服务端True,还是测试 grpc整个模块False
	server         proto.StreamServiceServer //grpc 服务端
	lstconn        []net.Conn                //客户端承载的 流 的汇总列表
}

func newmockClient(gserver proto.StreamServiceServer, size int) *mockClient {
	return &mockClient{
		server:  gserver,
		lstconn: make([]net.Conn, size),
	}
}

func (m *mockClient) SetServer(gserver proto.StreamServiceServer) {
	m.server = gserver
}

/****************************************/
/** 实现grpc StreamServiceClient 接口 ****/
/****************************************/

// StreamDual : 实现的grpc接口,获取 可以和服务端沟通的 流
// 将会生成 2个流，分别用于 客户端-服务端的单向流动 和 服务端-客户端的单向流动
func (m *mockClient) StreamDual(ctx context.Context,
	opts ...grpc.CallOption) (c proto.StreamService_StreamDualClient, err error) {

	chclient2server := make(chan proto.StreamBytes, 2)
	chserver2client := make(chan proto.StreamBytes, 2)
	servicestream := newserverstream(m, nil, chserver2client, chclient2server)
	var clientstream *clientstream
	if m.testgrpcServer {
		clientstream = newclientstream(m.server, servicestream, chclient2server, chserver2client, true)
	} else {
		clientstream = newclientstream(m.server, servicestream, chclient2server, chserver2client, false)
	}
	servicestream.Clientstream = clientstream

	c = clientstream
	return
}

// GetDomain : 实现 grpc接口,获取域名的ip
func (m *mockClient) GetDomain(ctx context.Context, req *proto.DomainReq,
	opts ...grpc.CallOption) (res *proto.IPRespose, err error) {
	return m.server.GetDomain(ctx, req)
}

func (m *mockClient) ResolvLocation(ctx context.Context, in *proto.LocationReq, opts ...grpc.CallOption) (location *proto.LocationResponse, err error) {
	location, err = m.server.ResolvLocation(ctx, in)
	return
}

/****************************************/
/** 以上完成grpc StreamServiceClient 接口 ****/
/****************************************/
