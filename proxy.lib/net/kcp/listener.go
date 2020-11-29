package kcp

import (
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/lancelotXie/proxy/proxy.lib/net/kcp/base"

	"github.com/pkg/errors"
	"github.com/xtaci/smux"
)

type listener struct {
	l              sync.Mutex
	base           net.Listener
	sessions       []*smux.Session
	sessionsClosed sync.WaitGroup
	newConns       chan net.Conn
	listening      bool
}

func newListenr(lis net.Listener) (l *listener) {
	l = new(listener)
	l.base = lis
	l.newConns = make(chan net.Conn)
	return
}

func listen(addr string) (lis *listener, err error) {
	conn, err := base.Listen(addr)
	if err != nil {
		return
	}

	l := newListenr(conn)
	lis = l
	return
}

func (l *listener) tryStartListen() {
	l.l.Lock()
	defer l.l.Unlock()
	if l.listening {
		return
	}
	l.listening = true

	go l.loopListen()
	go l.checkSessionsLoop()
}

func (l *listener) loopListen() {
	log.Println("开始复用监听：", l.Addr())
	for {
		conn, err := l.base.Accept()
		if err != nil {
			log.Println(err)
			break
		}

		log.Println("捕获 KCP 连接")
		go l.handle(conn)
	}
	log.Println("结束复用监听")
}

func (l *listener) handle(conn net.Conn) {
	config := smux.DefaultConfig()
	config.Version = 1
	config.MaxReceiveBuffer = 4194304
	config.MaxStreamBuffer = 2097152
	config.KeepAliveInterval = 10 * time.Second

	s, err := smux.Server(conn, config)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("成功建立 session 监听器")

	l.l.Lock()
	defer l.l.Unlock()
	l.sessions = append(l.sessions, s)
	l.sessionsClosed.Add(1)
	go l.handleSession(s)
}

func (l *listener) handleSession(s *smux.Session) {
	log.Println("开始 session 监听")
	for {
		stream, err := s.AcceptStream()
		if err != nil {
			log.Println(err)
			break
		}

		l.newConns <- stream
	}
	log.Println("结束 session 监听")
	l.sessionsClosed.Done()
}

// Accept 接收请求
func (l *listener) Accept() (conn net.Conn, err error) {
	l.tryStartListen()

	if !l.listening {
		err = errors.WithStack(io.EOF)
		return
	}

	conn, ok := <-l.newConns
	if !ok {
		err = errors.WithStack(io.EOF)
	} else {
		log.Println("成功接收复用 session")
	}
	return
}

func (l *listener) Close() (err error) {
	l.l.Lock()
	defer l.l.Unlock()
	if !l.listening {
		return
	}
	l.listening = false

	err = l.base.Close()

	for _, s := range l.sessions {
		s.Close()
	}

	l.sessionsClosed.Wait()
	close(l.newConns)
	for conn := range l.newConns {
		conn.Close()
	}

	return
}

func (l *listener) Addr() (addr net.Addr) {
	addr = l.base.Addr()
	return
}

func (l *listener) checkSessionsLoop() {
	for l.listening {
		l.checkSessions()
		time.Sleep(time.Minute)
	}
}

func (l *listener) checkSessions() {
	l.l.Lock()
	defer l.l.Unlock()

	aliveSessions := removeClosedSessions(l.sessions)
	l.sessions = aliveSessions
}

// removeClosedSessions 从集合中移除已关闭的 session
func removeClosedSessions(sessions []*smux.Session) (dst []*smux.Session) {
	for _, s := range sessions {
		if !s.IsClosed() {
			dst = append(dst, s)
		}
	}
	return
}
