package client

import (
	"context"
	"log"
	"net"
	"strconv"

	address "github.com/lancelotXie/proxy/proxy.lib/address"
	"github.com/lancelotXie/proxy/proxy.lib/ctxtransid"
	locationbase "github.com/lancelotXie/proxy/proxy.lib/location/base"
	"github.com/lancelotXie/proxy/proxy.lib/logger"
	proto "github.com/lancelotXie/proxy/proxy.server/grpc.server/proto"

	errors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

type streamclient struct {
	grpcclient proto.StreamServiceClient
}

func newstreamClient(serverip string, serverport int) (*streamclient, error) {
	client := &streamclient{}
	ipe := serverip + ":" + strconv.FormatInt(int64(serverport), 10)
	gconn, err := grpc.Dial(ipe, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client.grpcclient = proto.NewStreamServiceClient(gconn)
	return client, err
}

// GetConnByAddr :
func (c *streamclient) GetConnByAddr(remoteaddr net.Addr, localaddr net.Addr) (conn net.Conn, err error) {
	ip, port, network, ok := address.Unpack(remoteaddr)
	if !ok {
		return nil, errors.New("err: 远端地址转换出错")
	}
	ctx := ctxtransid.NewCTX(context.Background())
	ctx, cancelfunc := context.WithCancel(ctx)
	dualstream, err := c.grpcclient.StreamDual(ctx)
	log.Println("info:生成流，id：", ctxtransid.GetID(ctx))
	if err != nil {
		err = errors.WithStack(err)
		defer cancelfunc()
		return nil, err
	}
	req, err := proto.NewStreamBytesWithoutData(ip, port, network, ctxtransid.GetID(ctx))
	if err != nil {
		err = errors.WithStack(err)
		cancelfunc()
		return
	}
	dualstream.Send(req)
	conn = newclientConn(ctx, nil, remoteaddr, cancelfunc, dualstream)
	return
}

// Resolv :
func (c *streamclient) Resolv(domain string) (ip string, err error) {
	ctx := context.Background()
	req := &proto.DomainReq{
		Domain: domain,
	}
	res, err := c.grpcclient.GetDomain(ctx, req)
	if err != nil {
		return "", err
	}
	log.Println("info:解析dns:", req.Domain, res.IP)
	ip = res.IP
	return ip, nil
}

func dial(network string, ip string, port int) (conn net.Conn, err error) {
	streamServiceClient, err := defaultPool.GetClient()
	if err != nil {
		logger.Error(ctxtransid.NewCTX(context.Background()), err)
		return
	}
	ctx := ctxtransid.NewCTX(context.Background())
	ctx, cancelfunc := context.WithCancel(ctx)
	dualstream, err := streamServiceClient.StreamDual(ctx)
	if err != nil {
		logger.Error(ctx, err)
		cancelfunc()
		return
	}
	id := ctxtransid.GetID(ctx)
	remoteaddr := address.New(ip, port, network)
	req, err := proto.NewStreamBytesWithoutData(ip, port, network, id)
	if err != nil {
		cancelfunc()
		err = errors.WithStack(err)
		return
	}
	err = dualstream.Send(req)
	conn = newclientConn(ctx, nil, remoteaddr, cancelfunc, dualstream)
	logger.Info(ctx, "获得grpc组装的连接")
	return
}

func resolvLocation(ip string) (location locationbase.Location, err error) {
	client, err := defaultPool.GetClient()
	if err != nil {
		return
	}
	ctx := ctxtransid.NewCTX(context.Background())
	req := proto.NewLocationReq(ip)
	res, err := client.ResolvLocation(ctx, req)
	if err != nil {
		return
	}
	location = locationbase.Location(res.Location)
	return
}
