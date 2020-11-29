package client

import (
	"strconv"
	"testing"

	"github.com/lancelotXie/proxy/proxy.lib/address"
	"github.com/lancelotXie/proxy/proxy.lib/location"
	mockgrpc "github.com/lancelotXie/proxy/proxy.lib/mockgrpc"
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestResolv(t *testing.T) {
	var c GRPCClient
	var (
		ip            = "144.144.144.144"
		domain        = "www.baidu.com"
		invalidDomain = "www.google.com"
	)
	domainMap := make(map[string]string)
	domainMap[domain] = ip
	gclient, _ := mockgrpc.NewCS(1000, domainMap, nil)

	c = &streamclient{grpcclient: gclient}
	outip, err := c.Resolv(domain)
	assert.Nil(t, err)
	assert.Equal(t, ip, outip)

	ip2, err := c.Resolv(invalidDomain)
	assert.NotNil(t, err)
	assert.Empty(t, ip2)
}

func TestDial(t *testing.T) {
	var c GRPCClient
	var (
		ip      = "144.144.144.144"
		port    = 443
		network = "tcp"
		domain  = "www.baidu.com"
		msg1    = "client send"
	)
	domainMap := make(map[string]string)
	domainMap[domain] = ip
	gclient, _ := mockgrpc.NewCS(1000, domainMap, resFunc)

	c = &streamclient{grpcclient: gclient}

	conn, err := c.GetConnByAddr(address.New(ip, port, network), nil)
	assert.Nil(t, err)
	n, err := conn.Write([]byte(msg1))
	assert.Nil(t, err)
	assert.NotZero(t, n)

	data := make([]byte, bufferLen)
	len, err := conn.Read(data)
	assert.Nil(t, nil)
	len, err = conn.Read(data)
	assert.Nil(t, err)
	assert.Equal(t, msg1, string(data[:len]))

}

func resFunc(ip string, port int, net string) (res string) {
	res = ip + ":" + strconv.Itoa(int(port))
	return
}

func mocknewclient(serverip string, serverport int) (gconn *grpc.ClientConn, err error) {
	return &grpc.ClientConn{}, nil
}

func mockclose(conn *grpc.ClientConn) error { return nil }

func TestNewPool(t *testing.T) {
	var (
		size       = 10
		serverip   = "127.0.0.1"
		serverport = 8090
	)
	pool, err := newPool(size, serverip, serverport, mockNewClientConn)
	assert.Nil(t, err)
	assert.NotNil(t, pool)

	assert.Equal(t, size, pool.capacity)
	for _, v := range pool.lstconn {
		assert.NotNil(t, v)
	}
	assert.Equal(t, int32(0), pool.putnext)
}

func TestGetConn(t *testing.T) {
	var (
		size       = 10
		serverip   = "127.0.0.1"
		serverport = 8090
	)
	pool, err := newPool(size, serverip, serverport, mockNewClientConn)
	assert.Nil(t, err)
	for cur := 0; cur < pool.capacity*2; cur++ {
		mockclient, err := pool.GetClient()
		assert.Nil(t, err)
		assert.NotNil(t, mockclient)
		idx := cur % pool.capacity
		assert.Equal(t, pool.lstconn[idx], mockclient)
	}
}

func mockNewClientConn(string, int) (*grpc.ClientConn, error) {
	return &grpc.ClientConn{}, nil
}

func TestResolvLocation(t *testing.T) {
	var (
		size       = 1000
		serverip   = "127.0.0.1"
		serverport = 8090
		ip         = "144.144.144.144"
		loc        = "testlocation"
	)
	domainMap := make(map[string]string)
	domainMap[ip] = loc
	gclient, _ := mockgrpc.NewCS(size, domainMap, resFunc)

	pool, err := newPool(size, serverip, serverport, mockNewClientConn)
	assert.Nil(t, err)
	pool.lstconn[0] = gclient
	defaultPool = pool

	defaultLocationResolver = location.New(resolvLocation)

	res, err := ResolvLocation(ip)
	assert.NoError(t, err)
	assert.Equal(t, loc, res.String())
}

func TestStreamConn(t *testing.T) {
	var msgA = "DDDDDDDDDDD"
	var msgB = "EEEEEEEEEEE"
	cStream, sStream := mockgrpc.NewClientStreamAndServerStream()

	cStream.Send(proto.NewStreamBytes([]byte(msgA)))
	res, _ := sStream.Recv()
	assert.Equal(t, msgA, string(res.Data))
	sStream.Send(proto.NewStreamBytes([]byte(msgB)))
	res2, _ := cStream.Recv()
	assert.Equal(t, msgB, string(res2.Data))

}
