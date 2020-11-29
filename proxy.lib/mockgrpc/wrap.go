package mockgrpc

import (
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
)

// NewCS : 构造 grpc 客户端和服务端；
// 服务端“域名-IP”将会从domain map中取；
// 服务端 响应请求的 msg是由res生成的。
// connSize int: 指定了 客户端和服务端最大连接数(stream流)
// domain map: 包含了 服务端所有可以返回的 “域名-IP”
// res func: 指定 服务端 以何种方式处理 请求中数据/若无特殊要求，建议 使用nil
func NewCS(connSize int, domain map[string]string,
	res func(ip string, port int, network string) (res string)) (client proto.StreamServiceClient, server proto.StreamServiceServer) {
	mockclient := newmockClient(nil, connSize)
	mockserver := newmockServer(mockclient, connSize, domain, res)
	mockclient.SetServer(mockserver)
	return mockclient, mockserver
}

// NewClient :	测试 grpc server端时，将生成一个grpc serviceclient。可以由这个client发起 grpc请求
// 从而 测试 server端的代码。
func NewClient(gserver proto.StreamServiceServer, connSize int) (client proto.StreamServiceClient) {
	mockclient := newmockClient(gserver, connSize)
	mockclient.testgrpcServer = true
	client = mockclient
	return
}

// NewClientStreamAndServerStream : //构造2个 可以互相通信的 grpc stream
func NewClientStreamAndServerStream() (cStream proto.StreamService_StreamDualClient, sStream proto.StreamService_StreamDualServer) {
	client2server := make(chan proto.StreamBytes, 2)
	server2client := make(chan proto.StreamBytes, 2)

	clientStream := newclientstream(nil, nil, client2server, server2client, false)
	serverStream := newserverstream(nil, clientStream, server2client, client2server)
	clientStream.Serverstream = serverStream
	clientStream.sendDo = nil

	cStream = clientStream
	sStream = serverStream
	return
}
