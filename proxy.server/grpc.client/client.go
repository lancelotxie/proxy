package client

import (
	"strconv"

	"google.golang.org/grpc"
)

func newclient(serverip string, serverport int) (gconn *grpc.ClientConn, err error) {
	ipe := serverip + ":" + strconv.FormatInt(int64(serverport), 10)
	gconn, err = grpc.Dial(ipe, grpc.WithInsecure())
	return
}

func close(client *grpc.ClientConn) (err error) {
	err = client.Close()
	return
}
