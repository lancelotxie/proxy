package dial

import (
	"context"
	"net"
	"time"

	"github.com/lancelotXie/proxy/proxy.lib/logger"
	errors "github.com/pkg/errors"
)

// Dial : 连接到网络
type Dial interface {
	Conn2Web(context.Context, net.Addr) (net.Conn, error)
}

// WebServer : 实现Dial 接口
type WebServer struct {
}

// Conn2Web : 根据remote地址,发起tcp请求, 返回对应的 连接。默认超时30s
func (w *WebServer) Conn2Web(ctx context.Context, remote net.Addr) (conn net.Conn, err error) {
	if remote == nil {
		err = errors.New("err:remote 地址为nil。无法连接网络")
		logger.Error(ctx, err)
		return
	}
	ipe := remote.String()
	conn, err = net.DialTimeout("tcp", ipe, time.Second*30)
	if err != nil {
		logger.Error(ctx, "web 返回错误", err)
		return
	}
	logger.Info(ctx, "连接到web 成功")
	return
}
