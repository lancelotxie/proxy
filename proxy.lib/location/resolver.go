package location

import (
	"github.com/lancelotXie/proxy/proxy.lib/location/base"
	"github.com/lancelotXie/proxy/proxy.lib/resolver"
)

// ResolvFunc 解析地址的方法定义
type ResolvFunc func(ip string) (loc base.Location, err error)

// Resolver : 解析 IP 地址所在地
type Resolver interface {
	Resolv(ip string) (loc base.Location, err error)
}

type baseResolver struct {
	r resolver.StringResolver
}

func (r *baseResolver) Resolv(ip string) (loc base.Location, err error) {
	_loc, err := r.r.Resolv(ip)
	if err != nil {
		return
	}

	loc = base.Location(_loc)
	return
}

func newBaseResolver(base resolver.StringResolver) (r Resolver) {
	br := new(baseResolver)
	br.r = base
	r = br
	return
}
