/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-14 12:10:17
 * @LastEditTime: 2019-12-15 17:35:27
 */

package syncmap

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Marshal_UnMarshal(t *testing.T) {
	m0 := New()
	m1 := New()

	m0.Store("1", 1)
	m0.Store("2", 2)

	j, err := json.Marshal(m0)
	assert.NoError(t, err)

	err = json.Unmarshal(j, &m1)
	assert.NoError(t, err)

	v, ok := m1.Load("1")
	assert.Equal(t, 1, int(v.(float64)))
	assert.True(t, ok)

	v, ok = m1.Load("2")
	assert.Equal(t, 2, int(v.(float64)))
	assert.True(t, ok)
}

func Test_Recursive(t *testing.T) {
	m0 := New()
	m1 := New()
	m2 := New()

	m0.Store("1", m1)
	m1.Store("2", m2)
	m2.Store("3", "123")

	_m0 := New()

	j, err := json.Marshal(m0)
	assert.NoError(t, err)
	err = json.Unmarshal(j, &_m0)
	assert.NoError(t, err)

	v, ok := _m0.Load("1")
	assert.True(t, ok)
	_m1, ok := v.(Map)
	assert.True(t, ok)

	v, ok = _m1.Load("2")
	assert.True(t, ok)
	_m2, ok := v.(Map)
	assert.True(t, ok)

	v, ok = _m2.Load("3")
	assert.True(t, ok)
	s, ok := v.(string)
	assert.True(t, ok)
	assert.Equal(t, "123", s)
}
