package mockgrpc_test

import (
	"context"
	"testing"

	mockgrpc "github.com/lancelotXie/proxy/proxy.lib/mockgrpc"
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var (
		ip            = "144.144.144.144"
		domain        = "www.baidu.com"
		invalidDomain = "www.google.com"
	)
	domainMap := make(map[string]string)
	domainMap[domain] = ip
	c, _ := mockgrpc.NewCS(1000, domainMap, nil)

	req := &proto.DomainReq{
		Domain: domain,
	}
	resp, err := c.GetDomain(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, ip, resp.IP)

	req.Domain = invalidDomain
	resp, err = c.GetDomain(context.TODO(), req)
	assert.NotNil(t, err)
	assert.Nil(t, resp)
}

func TestDualStreamExternal(t *testing.T) {
	var (
		ip      = "144.144.144.144"
		port    = 443
		network = "tcp"
		domain  = "www.baidu.com"
	)
	domainMap := make(map[string]string)
	domainMap[domain] = ip
	c, _ := mockgrpc.NewCS(1000, domainMap, resFunc)

	req, err := proto.NewStreamBytesWithoutData(ip, port, network, "1")
	if !assert.NoError(t, err) {
		return
	}
	cStream, err := c.StreamDual(context.TODO())
	assert.Nil(t, err)
	err = cStream.Send(req)
	assert.Nil(t, err)

	response, err := cStream.Recv()
	assert.Nil(t, err)
	assert.Equal(t, string(req.Data), string(response.Data))
}

func resFunc(ip string, port int, network string) (res string) {
	return
}
