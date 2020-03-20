/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 16:38:12
 * @LastEditTime: 2019-12-07 16:51:17
 */

package renewconn

import (
	"bytes"
	"net"
)

type renewConn struct {
	net.Conn
	head *bytes.Buffer
}

func (c *renewConn) Read(b []byte) (n int, err error) {
	n, err = c.head.Read(b)
	if err == nil {
		return
	}

	// head 已经读完
	n, err = c.Conn.Read(b)
	return
}

func newRenewConn(ori net.Conn, head []byte) *renewConn {
	b := bytes.NewBuffer(head)
	conn := new(renewConn)
	conn.Conn = ori
	conn.head = b
	return conn
}
