package setting

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// ErrInvalidKind 非法的变量类型
var ErrInvalidKind = errors.New("invalid kind")

// kindBindSetting 带类型绑定的 setting
type kindBindSetting struct {
	setting
	types *sync.Map
}

func newkindBindSetting(base *setting) (s *kindBindSetting) {
	s = new(kindBindSetting)
	s.setting = *base
	s.types = new(sync.Map)
	return
}

// BindType 绑定类型
func (s *kindBindSetting) BindKind(key string, t reflect.Kind) {
	s.types.Store(key, t)
}

func (s *kindBindSetting) checkKind(key string, value interface{}) bool {
	validKind, ok := s.types.Load(key)
	if !ok {
		return true
	}

	t := reflect.TypeOf(value)
	kind := t.Kind()
	return validKind == kind
}

// Set 设置键值对，绑定过类型的 key 会检查类型是否匹配
func (s *kindBindSetting) Set(key string, value interface{}) (err error) {
	if !s.checkKind(key, value) {
		err = errors.WithStack(ErrInvalidKind)
		return
	}

	return s.setting.Set(key, value)
}

func (s *kindBindSetting) MarshalJSON() (b []byte, err error) {
	b, err = s.setting.MarshalJSON()
	return
}

func (s *kindBindSetting) UnmarshalJSON(b []byte) (err error) {
	err = s.setting.UnmarshalJSON(b)
	return
}
