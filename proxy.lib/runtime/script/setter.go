package script

import "log"

// Setter 设置器
type Setter interface {
	Get(key interface{}) (value interface{})
	Set(key interface{}, value interface{})
}

// BaseSetter 基础设置器
type BaseSetter struct {
	data map[interface{}]interface{}
}

// NewBaseSetter 构造新的 BaseSetter
func NewBaseSetter() (s *BaseSetter) {
	s = new(BaseSetter)
	s.data = make(map[interface{}]interface{})
	return
}

// Get 打印获取的值
func (bs *BaseSetter) Get(key interface{}) (value interface{}) {
	value = bs.data[key]
	log.Println("Get: ", key, " - ", value)
	return
}

// Set 打印设置的值
func (bs *BaseSetter) Set(key interface{}, value interface{}) {
	bs.data[key] = value
	log.Println("Set: ", key, " - ", value)
}
