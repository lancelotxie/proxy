package location

import (
	"github.com/lancelotXie/proxy/proxy.lib/location/api/ipapi"
	"github.com/lancelotXie/proxy/proxy.lib/location/base"
	resolver "github.com/lancelotXie/proxy/proxy.lib/resolver"
	"github.com/patrickmn/go-cache"
)

var defaultResolver = New(ipapi.Resolv)

// New : 返回Location的实例
func New(oriResolv ResolvFunc) (r Resolver) {
	resolv := func(in string) (out string, err error) {
		_out, err := oriResolv(in)
		if err != nil {
			return
		}
		out = _out.String()
		return
	}
	base := resolver.NewStringResolver(resolver.StringResolvFunc(resolv), cache.NoExpiration)
	r = newBaseResolver(base)

	return
}

// Resolv : 解析 IP 所在地
func Resolv(ip string) (loc base.Location, err error) {
	loc, err = defaultResolver.Resolv(ip)
	return
}
