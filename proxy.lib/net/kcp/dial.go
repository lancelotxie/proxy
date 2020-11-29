package kcp

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lancelotXie/proxy/proxy.lib/net/kcp/base"

	"github.com/xtaci/smux"
)

var defaultClientSession *smux.Session

var l4DefaultSession sync.Mutex

func buildSession(addr string) (err error) {
	conn, err := base.Dial(addr)
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	config := smux.DefaultConfig()
	config.Version = 1
	config.MaxReceiveBuffer = 4194304
	config.MaxStreamBuffer = 2097152
	config.KeepAliveInterval = 10 * time.Second

	if err = smux.VerifyConfig(config); err != nil {
		return
	}

	defaultClientSession, err = smux.Client(conn, config)
	if err == nil {
		log.Println("成功建立通向 ", addr, " 的 KCP 复用器")
	}
	return
}

func tryBuildSession(addr string) (err error) {
	l4DefaultSession.Lock()
	defer l4DefaultSession.Unlock()

	if defaultClientSession != nil {
		return
	}

	err = buildSession(addr)
	return
}

func reBuildSession(addr string) (err error) {
	l4DefaultSession.Lock()
	defer l4DefaultSession.Unlock()

	if defaultClientSession.IsClosed() {
		err = buildSession(addr)
	}
	return
}

// dial 拨号
func dial(addr string) (conn net.Conn, err error) {
	err = tryBuildSession(addr)
	if err != nil {
		return
	}

	conn, err = defaultClientSession.OpenStream()
	for err == io.ErrClosedPipe {
		err = reBuildSession(addr)
		if err != nil {
			break
		}

		conn, err = defaultClientSession.OpenStream()
		time.Sleep(time.Second)
	}
	if err == nil {
		log.Println("成功打开新的复用流")
	}
	return
}
