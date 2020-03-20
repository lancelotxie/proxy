package resolver

import (
	"time"

	"github.com/lancelot/proxy/proxy.lib/resolver/singleton"

	"github.com/lancelot/proxy/proxy.lib/resolver/base"
	"github.com/lancelot/proxy/proxy.lib/resolver/cache"
)

// newResolver 构造一个默认的解析器
func newResolver(resolv base.ResolvFunc, expiration time.Duration) (r base.Resolver) {
	br := base.New(resolv)

	// 结合 缓存层 与 单例层 实现解析请求合并
	cr := cache.New(br, expiration)
	sr := singleton.New(cr)

	r = sr
	return
}
