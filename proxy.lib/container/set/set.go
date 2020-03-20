/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-07 13:27:49
 * @LastEditTime: 2019-12-07 13:31:34
 */

package set

// Set 集合
type Set struct {
	data map[interface{}]struct{}
}

// NewSet 构造新集合
func NewSet() *Set {
	data := make(map[interface{}]struct{})
	s := new(Set)
	s.data = data
	return s
}

// Set 设置值
func (s *Set) Set(v interface{}) {
	s.data[v] = struct{}{}
}

// Exist 判断值是否存在
func (s Set) Exist(v interface{}) bool {
	_, ok := s.data[v]
	return ok
}
