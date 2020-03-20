/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-14 12:33:49
 * @LastEditTime: 2019-12-15 18:05:26
 */

package setting

import (
	"encoding/json"
	"strings"

	"github.com/lancelot/proxy/proxy.lib/container/syncmap"
)

const split = "."

// setting 配置管理器
type setting struct {
	m syncmap.Map
}

// newSetting 构造 setting
func newSetting() *setting {
	s := new(setting)
	s.m = syncmap.New()
	return s
}

func (s *setting) get(key string) (value interface{}, ok bool) {
	keysOfSubSettings, keyOfValue := splitKeys(key)
	subSetting, ok := s.getSubSettingRecursively(keysOfSubSettings, false)
	if !ok {
		return
	}

	value, ok = subSetting.getDirectly(keyOfValue)
	return
}

// Get 获取配置
func (s *setting) Get(key string) (value interface{}, ok bool) {
	return s.get(key)
}

func (s *setting) GetInt(key string) (value int, ok bool) {
	v, ok := s.get(key)
	if !ok {
		return
	}

	value, ok = toInt(v)
	return
}

func (s *setting) GetFloat(key string) (value float64, ok bool) {
	v, ok := s.get(key)
	if !ok {
		return
	}

	value, ok = v.(float64)
	return
}

func (s *setting) GetString(key string) (value string, ok bool) {
	v, ok := s.get(key)
	if !ok {
		return
	}

	value, ok = v.(string)
	return
}

func (s *setting) Set(key string, value interface{}) error {
	keysOfSubSettings, keyOfValue := splitKeys(key)
	subSetting, ok := s.getSubSettingRecursively(keysOfSubSettings, true)
	if !ok {
		panic("sub setting should be created")
	}

	subSetting.set(keyOfValue, value)
	return nil
}

// getDirectly 实际从底层获取 value
func (s *setting) getDirectly(key string) (value interface{}, ok bool) {
	value, ok = s.m.Load(key)
	return
}

// set 实际将 key/value 写入底层
func (s *setting) set(key string, value interface{}) {
	s.m.Store(key, value)
}

// getSubSetting 获取子 setting
func (s *setting) getSubSetting(key string) (sub *setting, ok bool) {
	value, ok := s.getDirectly(key)
	if !ok {
		return
	}
	sub, ok = value.(*setting)
	return
}

// getSubSettingRecursively 递归地获取子 setting，参数 force 表示是否强制创建不存在的子 setting
func (s *setting) getSubSettingRecursively(keys []string, force bool) (sub *setting, ok bool) {
	sub = s
	ok = true

	for _, key := range keys {
		// 迭代获取下一层 setting
		var sNext *setting
		sNext, ok = sub.getSubSetting(key)
		if ok {
			sub = sNext
			continue
		}

		// 子 setting 不命中，非强制模式，退出迭代
		if !force {
			break
		}

		// 构造子 setting
		sNext = newSetting()
		sub.set(key, sNext)
		sub = sNext
		ok = true
	}

	return
}

func (s *setting) MarshalJSON() (b []byte, err error) {
	b, err = json.Marshal(s.m)
	return
}

func (s *setting) UnmarshalJSON(b []byte) (err error) {
	err = json.Unmarshal(b, &s.m)

	s.m.Range(func(key, value interface{}) bool {
		_, ok := value.(syncmap.Map)
		if !ok {
			return true
		}

		// 子 Map 需要被转化为子 setting
		j, _ := json.Marshal(value)
		sub := newSetting()
		json.Unmarshal(j, &sub)
		s.m.Store(key, sub)
		return true
	})
	return
}

// splitKeys 拆分所有分代 setting 的 key 与 最终 value 的 key
func splitKeys(key string) (keysOfSubSettings []string, keyOfValue string) {
	keys := strings.Split(key, split)

	keysOfSubSettings = keys[:len(keys)-1]
	keyOfValue = keys[len(keys)-1]

	return
}

func toInt(src interface{}) (dst int, ok bool) {
	dst, ok = src.(int)
	if ok {
		return
	}

	// JSON 操作会将 int 转化为 float64
	_dst, ok := src.(float64)
	dst = int(_dst)
	return
}
