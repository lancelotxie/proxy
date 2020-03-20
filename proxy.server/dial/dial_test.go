package dial

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/lancelot/proxy/proxy.lib/address"

	"github.com/stretchr/testify/assert"
)

var req = "GET / HTTP/1.1\r\n" +
	"Host: www.huya.com\r\n" +
	"Connection: keep-alive\r\n" +
	"DNT: 1\r\n" +
	"Upgrade-Insecure-Requests: 1\r\n" +
	"User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36\r\n" +
	"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3" +
	"Accept-Encoding: gzip, deflate\r\n" +
	"Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7\r\n\r\n"

func TestConn2Web(t *testing.T) {
	var d Dial
	var webip = "182.131.6.243"
	var port int = 80
	d = &WebServer{}
	inaddr := address.New(webip, int(port), "")
	conn, err := d.Conn2Web(context.TODO(), inaddr)
	assert.Nil(t, err)
	assert.NotNil(t, conn)
	conn.Write([]byte(req))
	go func() {
		time.Sleep(time.Second * 2)
		os.Exit(0)
	}()
	for {
		rbytes := make([]byte, 1024)
		n, err := conn.Read(rbytes)
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		assert.Nil(t, err)
		assert.NotZero(t, n)
	}
}
