package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lancelot/proxy/proxy.lib/resolver/base"
	"github.com/patrickmn/go-cache"
)

func Test_Resolv(t *testing.T) {
	resolv := func(in string) (out interface{}, err error) {
		time.Sleep(time.Millisecond * 100)
		return
	}
	r := base.New(resolv)
	cr := New(r, cache.NoExpiration)

	t0 := time.Now()
	cr.Resolv("1")
	t1 := time.Now()
	cr.Resolv("1")
	t2 := time.Now()

	d0 := t1.Sub(t0)
	d1 := t2.Sub(t1)

	assert.GreaterOrEqual(t, d0.Milliseconds(), int64(100))
	assert.Greater(t, d0.Milliseconds(), d1.Milliseconds())
}
