package server

import (
	"context"
	"net"
	"testing"

	"github.com/lancelotXie/proxy/proxy.lib/address"
	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	"github.com/lancelotXie/proxy/proxy.lib/mockgrpc"
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"

	errors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var serverport int = 8090

func TestAccept(t *testing.T) {
	var (
		ip          = "144.144.144.144"
		port    int = 443
		network     = "tcp"
		domain      = "www.baidu.com"
		count       = 1000
	)
	domainMap := make(map[string]string)
	domainMap[domain] = ip
	s := newgrpcServer(&net.TCPListener{})
	gclient := mockgrpc.NewClient(s, count)

	ctx := ctxtransid.NewCTX(context.Background())
	clientstream, err := gclient.StreamDual(ctx)
	assert.Nil(t, err)
	inaddr := address.New(ip, int(port), network)
	req, err := proto.NewStreamBytesWithoutData(ip, int(port), network, ctxtransid.GetID(ctx))
	if err != nil {
		err = errors.WithStack(err)
		t.FailNow()
	}
	err = clientstream.Send(req)
	assert.NoError(t, err)
	_, outconn, err := s.Accept()

	assert.Nil(t, err)
	assert.Equal(t, inaddr.String(), outconn.RemoteAddr().String())
	assert.Equal(t, inaddr.Network(), outconn.RemoteAddr().Network())
	outconn.Close()
}

func TestWrite(t *testing.T) {
	var (
		inip        = "144.144.144.144"
		inport  int = 443
		network     = "tcp"
		domain      = "www.baidu.com"
		count       = 1000
		n           = 0
	)
	domainMap := make(map[string]string)
	domainMap[domain] = inip
	s := newgrpcServer(&net.TCPListener{})
	ctx := ctxtransid.NewCTX(context.Background())
	gclient := mockgrpc.NewClient(s, count)
	clientStream, err := gclient.StreamDual(ctx)
	assert.Nil(t, err)
	req, err := proto.NewStreamBytesWithoutData(inip, int(inport), network, ctxtransid.GetID(ctx))
	if err != nil {
		err = errors.WithStack(err)
		t.FailNow()
	}
	err = clientStream.Send(req)
	assert.Nil(t, err)

	_, serverconn, err := s.Accept()
	assert.Nil(t, err)

	cwbytes := []byte("client write")
	req = proto.NewStreamBytes(cwbytes)
	err = clientStream.Send(req)
	assert.Nil(t, err)

	srbytes := make([]byte, 1024)
	n, err = serverconn.Read(srbytes)
	assert.Nil(t, err)
	assert.Equal(t, string(cwbytes), string(srbytes[:n]))

	swstring := "server write"
	n, err = serverconn.Write([]byte(swstring))
	assert.NotZero(t, n)
	assert.Nil(t, err)

	resp, err := clientStream.Recv()
	assert.Nil(t, err)
	creadstring := string(resp.Data)
	assert.Nil(t, err)
	assert.Equal(t, string(swstring), creadstring)
}

func TestGetDomain(t *testing.T) {
	var (
		ip     = "144.144.144.144"
		domain = "www.baidu.com"
		count  = 1000
	)
	domainMap := make(map[string]string)
	domainMap[domain] = ip
	s := newgrpcServer(&net.TCPListener{})
	s.resolver = &mockResolver{domainMap: domainMap}
	gclient := mockgrpc.NewClient(s, count)

	req := &proto.DomainReq{
		Domain: domain,
	}
	res, err := gclient.GetDomain(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, ip, res.IP)
}

type mockResolver struct{ domainMap map[string]string }

func (m *mockResolver) Resolv(domain string) (ip string, err error) {
	ip, ok := m.domainMap[domain]
	if !ok {
		err = errors.New("cannot found domain")
	}
	return
}

func TestResolvLocation(t *testing.T) {
	var (
		ip       = "144.144.144.144"
		location = "Others"
		count    = 1000
	)
	s := newgrpcServer(&net.TCPListener{})
	gclient := mockgrpc.NewClient(s, count)

	ctx := context.TODO()
	req := proto.NewLocationReq(ip)
	res, err := gclient.ResolvLocation(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, location, res.Location)
}
