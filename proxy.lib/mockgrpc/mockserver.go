package mockgrpc

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
)

type streamServiceServer interface {
	proto.StreamServiceServer
}

// mockServer : 模拟 grpc 服务端
type mockServer struct {
	streamLocker *sync.Mutex
	curStream    int                                                    //当前正在处理的流的下标
	client       proto.StreamServiceClient                              //相匹配的grpc 客户端
	lstStream    []proto.StreamService_StreamDualServer                 //当前服务端可以处理的流(包含nil)
	mapDomain    map[string]string                                      //服务端 可以处理的 域名-IP 字典,IP-所在地 字典
	resFunc      func(ip string, port int, network string) (res string) //处理请求数据的方式
}

func newmockServer(gclient proto.StreamServiceClient, size int,
	mapDomain map[string]string,
	resfunc func(ip string, port int, network string) (res string)) *mockServer {
	return &mockServer{
		streamLocker: &sync.Mutex{},
		curStream:    0,
		client:       gclient,
		lstStream:    make([]proto.StreamService_StreamDualServer, size),
		mapDomain:    mapDomain,
		resFunc:      resfunc,
	}
}

/******************/
/*   实现grpc接口  */
/******************/

// StreamDual :实现的grpc 接口,每处理一个请求，将当前流 计入 服务端 的流列表
func (m *mockServer) StreamDual(s proto.StreamService_StreamDualServer) (err error) {
	m.streamLocker.Lock()
	defer m.streamLocker.Unlock()
	m.lstStream[m.curStream] = s
	m.curStream++
	go m.respondNow(s)
	return nil
}

// GetDomain :实现的grpc接口
func (m *mockServer) GetDomain(ctx context.Context, req *proto.DomainReq) (res *proto.IPRespose, err error) {
	domain := req.Domain
	ip, ok := m.mapDomain[domain]
	if ok {
		res = &proto.IPRespose{IP: ip}
		return
	}
	err = errors.New("err:未查找到" + domain + "对应的ip")
	return
}

// responseNow : 回复请求
func (m *mockServer) respondNow(stream proto.StreamService_StreamDualServer) {
	for {
		req, err := stream.Recv()
		if err != nil {
			fmt.Println(errors.WithStack(err))
			return
		}
		if m.resFunc == nil {
			fmt.Println("server won't respond,because resFunc is nil")
			return
		}

		res := &proto.StreamBytes{
			Data: []byte(req.Data),
		}
		stream.Send(res)
	}
}

func (m *mockServer) ResolvLocation(ctx context.Context, req *proto.LocationReq) (location *proto.LocationResponse, err error) {
	loc := m.mapDomain[req.IP]
	location = proto.NewLocationResponse(loc)
	return
}

/******************/
/*   完成grpc接口  */
/******************/
