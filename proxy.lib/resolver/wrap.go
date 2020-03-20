package resolver

import (
	"time"

	"github.com/lancelot/proxy/proxy.lib/resolver/base"
	stringresolver "github.com/lancelot/proxy/proxy.lib/resolver/string"
)

// ResolvFunc 解析函数，即 base.ResolvFunc
type ResolvFunc = base.ResolvFunc

// StringResolvFunc string 输出的解析函数，即 stringresolver.ResolvFunc
type StringResolvFunc = stringresolver.ResolvFunc

// Resolver 解析器，即 base.Resolver
type Resolver = base.Resolver

// StringResolver 输出为 string 的解析器，即 stringresolver.Resolver
type StringResolver = stringresolver.Resolver

// New 构造一个解析器，expiration 表示缓存超时的时间
func New(resolv ResolvFunc, expiration time.Duration) (r base.Resolver) {
	r = newResolver(resolv, expiration)
	return
}

// NewStringResolver 构造一个只允许 string 输出的解析器， expiration 表示缓存超时的时间
func NewStringResolver(resolv StringResolvFunc, expiration time.Duration) (r stringresolver.Resolver) {
	_resolv := func(in string) (out interface{}, err error) {
		_out, err := resolv(in)
		// 将 string out 转化为 interface{} out
		out = _out
		return
	}

	_r := newResolver(_resolv, expiration)
	r = stringresolver.New(_r)
	return
}
