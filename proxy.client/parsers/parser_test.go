/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 13:51:55
 * @LastEditTime: 2019-12-07 14:06:47
 */

package parsers

import (
	"testing"

	"github.com/eaglexiang/go/tunnel"
	"github.com/stretchr/testify/assert"
)

func Test_Register(t *testing.T) {
	err := RegisterParser(&mockParser{})
	assert.NoError(t, err)

	err = RegisterParser(&mockParser{})
	assert.Error(t, err)
}

func Test_GetParser(t *testing.T) {
	mockConn := tunnel.NewVirtualConn()

	RegisterParser(&mockParser{})

	mockConn.Write([]byte("valid"))
	p, _, err := GetParser(mockConn)
	assert.NoError(t, err)
	assert.Equal(t, "mock", p.String())

	mockConn.Write([]byte("AAA"))
	_, _, err = GetParser(mockConn)
	assert.Error(t, err)
}
