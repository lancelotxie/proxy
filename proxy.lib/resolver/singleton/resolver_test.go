package singleton

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lancelot/proxy/proxy.lib/resolver/base"
)

func Test_singletonResolver(t *testing.T) {
	count1 := 0
	count2 := 0
	resolv := func(in string) (out interface{}, err error) {
		if in == "1" {
			count1++
		} else if in == "2" {
			count2++
			count2++
			count2++
		}
		return
	}
	br := base.New(resolv)
	sr := New(br)

	wg := sync.WaitGroup{}
	tryCount := 10000
	wg.Add(tryCount)
	for i := 0; i < tryCount; i++ {
		go func() {
			sr.Resolv("1")
			sr.Resolv("2")
			wg.Done()
		}()
	}
	wg.Wait()

	// 如果单例生效，那么所有相同 in 的 count 的 ++ 行为应该保持线性
	// ++ 行为不应该出现线程安全问题
	assert.Equal(t, tryCount, count1)
	assert.Equal(t, tryCount*3, count2)
}
