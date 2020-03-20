package singleton

import (
	"sync"

	"github.com/lancelot/proxy/proxy.lib/resolver/base"
)

type doFunc func()

// singletonResolver 单例解析器
type singletonResolver struct {
	base.Resolver
	requsts    map[string]chan struct{}
	l4Requests *sync.Mutex
}

// New 构造一个单例解析器，所有相同 in 的解析操作都会单例执行
func New(r base.Resolver) (mr base.Resolver) {
	_mr := &singletonResolver{
		Resolver:   r,
		requsts:    make(map[string]chan struct{}),
		l4Requests: new(sync.Mutex),
	}
	mr = _mr
	return
}

func (r *singletonResolver) Resolv(in string) (out interface{}, err error) {
	done, ok, do := r.tryRegisterRequest(in)
	for ok {
		<-done
		done, ok, do = r.tryRegisterRequest(in)
	}
	defer do()

	out, err = r.Resolver.Resolv(in)
	return
}

// tryRegisterRequest 尝试注册单例请求。ok 为是否存在既有注册，当存在，应等待 done 信号，当不存在，应在完成
// 事务后执行 do 回调
func (r *singletonResolver) tryRegisterRequest(req string) (done chan struct{}, ok bool, do doFunc) {
	r.l4Requests.Lock()
	defer r.l4Requests.Unlock()

	// 查询是否已有注册
	done, ok = r.requsts[req]
	if ok {
		return
	}

	// 注册请求
	done = make(chan struct{})
	r.requsts[req] = done

	// 制造 done 信号，注销请求
	do = func() {
		r.l4Requests.Lock()
		defer r.l4Requests.Unlock()

		delete(r.requsts, req)
		close(done)
	}
	return
}
