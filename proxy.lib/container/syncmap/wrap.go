/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-14 12:15:31
 * @LastEditTime: 2019-12-14 12:16:15
 */

package syncmap

// New 构造新的并发安全 Map
func New() Map {
	return &syncmap{}
}
