package dns

import (
	"time"

	resolver "github.com/lancelot/proxy/proxy.lib/resolver"

	mynet "github.com/eaglexiang/go/net"
	"github.com/pkg/errors"
)

// ResolvFunc DNS 解析的函数
type ResolvFunc func(domain string) (ip string, err error)

// Resolver DNS 解析器
type Resolver interface {
	Resolv(domain string) (ip string, err error)
}

type dnsResolver struct {
	resolver.StringResolver
}

func (r *dnsResolver) Resolv(domain string) (ip string, err error) {
	ip, err = r.StringResolver.Resolv(domain)
	return
}

func newDNSResolver(resolv ResolvFunc) (r Resolver) {
	base := resolver.NewStringResolver(resolver.StringResolvFunc(resolv), time.Hour*2)
	_r := &dnsResolver{
		StringResolver: base,
	}
	r = _r
	return
}

func defaultResolv(domain string) (ip string, err error) {
	t := mynet.TypeOfAddr(domain)
	switch t {
	case mynet.DomainAddr:
		ip, err = mynet.ResolvIPv4(domain)
	case mynet.IPv4Addr, mynet.IPv6Addr:
		ip = domain
	default:
		err = errors.New("invalid address type")
	}
	return
}
