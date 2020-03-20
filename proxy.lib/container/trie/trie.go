package trie

import (
	"github.com/lancelot/proxy/proxy.lib/container/array"
)

// Trie 字典树
type Trie struct {
	root *trieNode
}

// New 构造新的字典树
func New() (t *Trie) {
	root := newNode()

	t = new(Trie)
	t.root = root
	return
}

// Set 将 a 填入 trie 中，如果成功则返回 true，如果失败（该串已存在）则返回 false
func (t *Trie) Set(a array.Array) (ok bool) {
	l := a.Len()

	node := t.root

	for i := 0; i < l; i++ {
		ele := a.Get(i)
		node = node.forceGetChild(ele)
	}

	if node.colored {
		return
	}

	ok = true
	node.colored = true
	return
}

// Exist 判断 trie 中是否存在 a
func (t *Trie) Exist(a array.Array) (ok bool) {
	l := a.Len()

	node := t.root

	for i := 0; i < l; i++ {
		ele := a.Get(i)
		node, ok = node.find(ele)
		if !ok {
			break
		}
	}
	if !ok {
		return
	}

	ok = node.colored

	return
}

// Match 判断 trie 中是否有 a 的头部子串
func (t *Trie) Match(a array.Array) (ok bool) {
	l := a.Len()

	node := t.root

	for i := 0; i < l; i++ {
		ele := a.Get(i)
		var found bool
		node, found = node.find(ele)
		if !found {
			break
		}

		if node.colored {
			ok = true
			break
		}
	}

	return
}
