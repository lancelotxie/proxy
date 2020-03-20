/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-15 17:53:41
 * @LastEditTime: 2019-12-15 18:06:48
 */

package setting

import (
	"reflect"
)

// Setting 配置管理器
type Setting interface {
	Get(key string) (value interface{}, ok bool)
	GetInt(key string) (value int, ok bool)
	GetFloat(key string) (value float64, ok bool)
	GetString(key string) (value string, ok bool)
	Set(key string, value interface{}) error
	BindKind(key string, kind reflect.Kind)
	Hook(key string, h Hooker)
}

// New 构造新的配置管理器
func New() Setting {
	s := newSetting()
	ks := newkindBindSetting(s)
	hs := newhookSetting(ks)
	return hs
}

var s = New()

// Hook 注册钩子
func Hook(key string, h Hooker) {
	s.Hook(key, h)
}

// Get 获取配置
func Get(key string) (value interface{}, ok bool) {
	return s.Get(key)
}

// Set 设置配置
func Set(key string, value interface{}) {
	s.Set(key, value)
}

// GetInt 获取 int 型配置
func GetInt(key string) (value int, ok bool) {
	return s.GetInt(key)
}

// GetFloat 获取 float64 型配置
func GetFloat(key string) (value float64, ok bool) {
	return s.GetFloat(key)
}

// GetString 获取 string 型配置
func GetString(key string) (value string, ok bool) {
	return s.GetString(key)
}
