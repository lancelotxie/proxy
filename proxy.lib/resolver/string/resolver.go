package string

import (
	"github.com/lancelotXie/proxy/proxy.lib/resolver/base"

	"github.com/pkg/errors"
)

// ErrOutputNotString 输出并不是 string 类型
var ErrOutputNotString = errors.New("type of output is not string")

// ResolvFunc 进出皆为 string 的解析函数
type ResolvFunc func(in string) (out string, err error)

// Resolver 进出皆为 string 的解析器
type Resolver interface {
	Resolv(in string) (out string, err error)
}

// New 构造一个只允许 string 类型输出的解析器
func New(r base.Resolver) (sr Resolver) {
	_sr := &resolver{
		Resolver: r,
	}
	sr = _sr
	return
}

type resolver struct {
	base.Resolver
}

func (r *resolver) Resolv(in string) (out string, err error) {
	_out, err := r.Resolver.Resolv(in)
	if err != nil {
		return
	}

	out, ok := _out.(string)
	if !ok {
		err = errors.WithStack(ErrOutputNotString)
	}
	return
}
