package flow

import (
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func sendwrite(leftconn net.Conn, rightconn net.Conn,
	leftwrite string, rightwrite string, lclone net.Conn, rclone net.Conn,
	waitg *sync.WaitGroup, t *testing.T) {
	waitg.Add(4)
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			wlen, err := leftconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			wlen, err := rightconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := rightconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := leftconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()
	waitg.Wait()

	leftconn.Close()
	rightconn.Close()
	lclone.Close()
	rclone.Close()
}

func TestMockconn(t *testing.T) {
	m := &mockconn{}
	var leftconn net.Conn
	var rightconn net.Conn
	var leftwrite = "left write"
	var rightwrite = "right write"
	leftconn, rightconn = m.GetConn()
	waitg := &sync.WaitGroup{}
	waitg.Add(4)
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			wlen, err := leftconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()
	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			wlen, err := rightconn.Write([]byte(cwrite))
			assert.NotZero(t, wlen)
			assert.Nil(t, err)
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := leftwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := rightconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()

	go func() {
		for num := 0; num < 4; num++ {
			cwrite := rightwrite + strconv.Itoa(num)
			crbytes := make([]byte, 1024)
			wlen, err := leftconn.Read(crbytes)
			if err == nil {
				assert.NotZero(t, wlen)
				assert.Equal(t, cwrite, string(crbytes[:wlen]))
			} else if err != io.EOF {
				log.Println("client read err:", err)
			}
		}
		waitg.Done()
	}()
	waitg.Wait()

	leftconn.Close()
	rightconn.Close()
}
