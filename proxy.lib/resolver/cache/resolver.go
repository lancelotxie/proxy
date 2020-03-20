package cache

import (
	"time"

	"github.com/lancelot/proxy/proxy.lib/resolver/base"

	"github.com/patrickmn/go-cache"
)

type cacheResolver struct {
	base.Resolver
	c *cache.Cache
}

func (r *cacheResolver) Resolv(in string) (out interface{}, err error) {
	out, ok := r.c.Get(in)
	if ok {
		return
	}

	out, err = r.Resolver.Resolv(in)
	if err != nil {
		return
	}

	r.c.SetDefault(in, out)
	return
}

// New 构造一个带缓存的解析器
func New(r base.Resolver, expiration time.Duration) (cr base.Resolver) {
	cleanupInterval := expiration / 2
	c := cache.New(expiration, cleanupInterval)

	_cr := &cacheResolver{
		Resolver: r,
		c:        c,
	}
	cr = _cr

	return
}
