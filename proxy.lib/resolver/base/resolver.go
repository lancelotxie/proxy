package base

// ResolvFunc 通用的解析函数声明
type ResolvFunc func(string) (interface{}, error)

// Resolver 通用的解析器声明
type Resolver interface {
	Resolv(string) (interface{}, error)
}

// baseResolver 最基本的解析器实现
type baseResolver ResolvFunc

// Resolv 解析 in，得到 out
func (br baseResolver) Resolv(in string) (out interface{}, err error) {
	out, err = br(in)
	return
}

// New 构造一个基本的解析器
func New(resolv ResolvFunc) (r Resolver) {
	r = baseResolver(resolv)
	return
}
