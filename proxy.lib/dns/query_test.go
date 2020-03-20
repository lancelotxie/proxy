/*
 * @Author: EagleXiang
 * @LastEditors: EagleXiang
 * @Email: eagle.xiang@outlook.com
 * @Github: https://github.com/eaglexiang
 * @Date: 2019-12-06 22:12:31
 * @LastEditTime: 2019-12-17 23:04:58
 */

package dns

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Resolver_Resolv(t *testing.T) {
	r := New(defaultResolv)
	ep0 := time.Now()
	ip, err := r.Resolv("www.baidu.com")
	ep1 := time.Now()
	assert.NoError(t, err)
	t.Log(ip)

	ep2 := time.Now()
	ip, err = r.Resolv("www.baidu.com")
	ep3 := time.Now()
	assert.NoError(t, err)
	t.Log(ip)

	// 第二次请求触发缓存，应该比第一次快
	d0 := ep1.Sub(ep0)
	d1 := ep3.Sub(ep2)
	f := func() bool {
		return d0 > d1
	}
	assert.Condition(t, f)
	t.Log("d0: ", d0, " d1: ", d1)
}
