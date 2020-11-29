package configuration

import (
	"os"
	"reflect"

	"github.com/lancelotXie/proxy/proxy.lib/path"
	"github.com/lancelotXie/proxy/proxy.lib/setting"

	"github.com/pkg/errors"
)

// ErrRemoteAddrNotFound 远端地址不存在
var ErrRemoteAddrNotFound = errors.New("remote addr not found")

var c = newConfiguration()

// Set 设置 key/value，每次写入成功都会覆写配置文件
func Set(key string, value interface{}) (err error) {
	err = c.Set(key, value)
	return
}

// Get 获取 Value
func Get(key string) (value interface{}, ok bool) {
	return c.Get(key)
}

// GetString 获取 string 类型的 value
func GetString(key string) (value string, ok bool) {
	return c.GetString(key)
}

// GetInt 获取 int 类型的 value
func GetInt(key string) (value int, ok bool) {
	return c.GetInt(key)
}

// GetFloat 获取 float 类型的 value
func GetFloat(key string) (value float64, ok bool) {
	return c.GetFloat(key)
}

// BindKind 绑定 key 的类型
func BindKind(key string, kind reflect.Kind) {
	c.BindKind(key, kind)
}

// Hook 注册钩子
func Hook(key string, h setting.Hooker) {
	c.Hook(key, h)
}

// Save 保存配置到文件
func Save() (err error) {
	mkdir()

	fileName, err := path.ConfigFile()
	if err != nil {
		return
	}
	err = c.Save(fileName)
	return
}

// Load 从文件读取配置
func Load() (err error) {
	fileName, err := path.ConfigFile()
	if err != nil {
		return
	}
	err = c.Load(fileName)
	return
}

// Clear 清空配置
func Clear() {
	c = newConfiguration()
}

func mkdir() {
	dirName, err := path.ConfigDir()
	if err != nil {
		return
	}

	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
}
