/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-18 21:22:52
 * @LastEditTime: 2019-12-18 21:41:31
 */

package path

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Home(t *testing.T) {
	home, err := Home()
	assert.NoError(t, err)
	t.Log(home)
}

func Test_Init(t *testing.T) {
	_, err := ConfigFile()
	assert.Error(t, err)

	Init(Server)
	path, err := ConfigFile()
	assert.NoError(t, err)
	t.Log(path)
}
