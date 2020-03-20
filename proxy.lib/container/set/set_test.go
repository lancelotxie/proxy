/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 13:31:41
 * @LastEditTime: 2019-12-07 13:34:04
 */

package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Set(t *testing.T) {
	s := NewSet()

	r0 := s.Exist("111")
	r1 := s.Exist("222")
	assert.False(t, r0)
	assert.False(t, r1)

	s.Set("111")

	r0 = s.Exist("111")
	r1 = s.Exist("222")
	assert.True(t, r0)
	assert.False(t, r1)
}
