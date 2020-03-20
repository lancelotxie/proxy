package base

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_kcp(t *testing.T) {
	lis, err := Listen("127.0.0.1:8877")
	if !assert.NoError(t, err) {
		return
	}
	defer lis.Close()
	log.Println("开始监听")

	done := make(chan struct{})
	go func() {
		conn, err := lis.Accept()
		if !assert.NoError(t, err) {
			return
		}
		log.Println("建立连接")

		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "test msg", string(buf[:n]))

		close(done)
	}()

	conn, err := Dial("127.0.0.1:8877")
	if !assert.NoError(t, err) {
		return
	}
	log.Println("拨号成功")

	buf := bytes.NewBufferString("test msg")
	_, err = buf.WriteTo(conn)
	assert.NoError(t, err)
	conn.Close()

	<-done
}
