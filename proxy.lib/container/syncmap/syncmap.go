/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-14 12:02:16
 * @LastEditTime: 2019-12-14 18:21:54
 */

package syncmap

import (
	"encoding/json"
	"sync"
)

// Map 映射表
type Map interface {
	Delete(key interface{})
	Load(key interface{}) (value interface{}, ok bool)
	LoadOrStore(key interface{}, value interface{}) (actual interface{}, loaded bool)
	Range(f func(key interface{}, value interface{}) bool)
	Store(key interface{}, value interface{})
	MarshalJSON() (b []byte, err error)
	UnmarshalJSON(b []byte) (err error)
}

type syncmap struct {
	sync.Map
}

// MarshalJSON 实现 json.Marshaler 接口
func (sm *syncmap) MarshalJSON() (b []byte, err error) {
	m := syncmap2map(&sm.Map)
	b, err = json.Marshal(m)
	return
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (sm *syncmap) UnmarshalJSON(b []byte) (err error) {
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}

	map2syncmap(m, &sm.Map)
	return
}

func syncmap2map(sm *sync.Map) (m map[string]interface{}) {
	m = make(map[string]interface{})

	sm.Range(func(key, value interface{}) bool {
		_key, ok := key.(string)
		if !ok {
			return true
		}
		m[_key] = value
		return true
	})

	return m
}

func map2syncmap(m map[string]interface{}, sm *sync.Map) {
	for key, value := range m {
		_, ok := value.(map[string]interface{})
		if !ok {
			sm.Store(key, value)
			continue
		}

		// value 是子 map
		j, _ := json.Marshal(value)
		sub := New()
		json.Unmarshal(j, &sub)
		sm.Store(key, sub)
	}
}
