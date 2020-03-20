/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-14 16:39:38
 * @LastEditTime: 2019-12-15 17:52:25
 */

package setting

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_splitKeys(t *testing.T) {
	key := "1.2.3"
	keysOfSubs, keyOfValue := splitKeys(key)
	assert.Equal(t, []string{"1", "2"}, keysOfSubs)
	assert.Equal(t, "3", keyOfValue)
}

func Test_getSubSettingRecursively(t *testing.T) {
	m := newSetting()

	keys := []string{"1", "2", "3"}

	// 初次获取，子 setting 不存在
	_, ok := m.getSubSettingRecursively(keys, false)
	assert.False(t, ok)

	// 再次获取，子 setting 仍然不存在
	_, ok = m.getSubSettingRecursively(keys, false)
	assert.False(t, ok)

	// 强制模式，子 setting 会被创建
	_, ok = m.getSubSettingRecursively(keys, true)
	assert.True(t, ok)

	// 子 setting 已经被创建过
	_, ok = m.getSubSettingRecursively(keys, false)
	assert.True(t, ok)
}

func Test_SetGet(t *testing.T) {
	m := New()

	_, ok := m.Get("123")
	assert.False(t, ok)

	m.Set("123", 123)
	v, ok := m.Get("123")
	assert.True(t, ok)
	assert.Equal(t, 123, v)

	m.Set("3.2.1", 321)
	v, ok = m.Get("3.2.1")
	assert.True(t, ok)
	assert.Equal(t, 321, v)
}

func Test_MarshalJSON(t *testing.T) {
	s := newSetting()

	_3 := newSetting()

	s.Set("1", 1)
	s.Set("2", "2")
	s.Set("3", _3)

	_3.Set("4", "444")

	j, err := json.Marshal(s)
	assert.NoError(t, err)

	str := string(j)
	assert.Equal(t, `{"1":1,"2":"2","3":{"4":"444"}}`, str)
}

func Test_UnmarshalJSON(t *testing.T) {
	str := `{"1":1,"2":"2","3":{"4":"444"}}`

	s := newSetting()
	err := json.Unmarshal([]byte(str), &s)
	assert.NoError(t, err)

	v, ok := s.Get("1")
	assert.True(t, ok)
	assert.Equal(t, float64(1), v)

	v, ok = s.Get("2")
	assert.True(t, ok)
	assert.Equal(t, "2", v)

	v, ok = s.Get("3")
	assert.True(t, ok)
	v3, ok := v.(*setting)
	assert.True(t, ok)

	v, ok = v3.Get("4")
	assert.True(t, ok)
	assert.Equal(t, "444", v)
}

func Test_setting_JSON(t *testing.T) {
	m0 := newSetting()

	m0.Set("123", 123)
	m0.Set("3.2.1", 321)

	m1 := newSetting()

	b, err := json.Marshal(m0)
	assert.NoError(t, err)

	err = json.Unmarshal(b, &m1)
	assert.NoError(t, err)

	v0, ok := m1.Get("123")
	assert.True(t, ok)
	assert.Equal(t, 123, int(v0.(float64)))

	v1, ok := m1.Get("3.2.1")
	assert.True(t, ok)

	assert.Equal(t, 321, int(v1.(float64)))
}
