/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 17:10:33
 * @LastEditTime: 2019-12-07 17:17:30
 */

package renewconn

import (
	"io/ioutil"
	"testing"

	"github.com/eaglexiang/go/tunnel"
	"github.com/stretchr/testify/assert"
)

func Test_Read(t *testing.T) {
	mock := tunnel.NewVirtualConn()
	m0 := "msg0"
	m1 := "msg1"
	m2 := "msg2"

	mock.Write([]byte(m0))
	renew := Renew(mock, []byte(m1))
	mock.Write([]byte(m2))
	mock.Close()

	b, err := ioutil.ReadAll(renew)
	assert.NoError(t, err)
	r := string(b)
	assert.Equal(t, m1+m0+m2, r, r)
}
