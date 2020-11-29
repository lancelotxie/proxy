package trie

import "github.com/lancelotXie/proxy/proxy.lib/container/array"

// InvertedTrie 逆序读写 array 的字典树
type InvertedTrie struct {
	Trie
}

// NewInverted 构造 InvertedTrie
func NewInverted() (t *InvertedTrie) {
	t = new(InvertedTrie)
	t.Trie = *New()
	return
}

// Set 将 a 逆序填入 trie 中
func (it *InvertedTrie) Set(a array.Array) (ok bool) {
	a = a.Inverted()
	ok = it.Trie.Set(a)
	return
}

// Exist 判断 trie 中是否存在逆序的 a
func (it *InvertedTrie) Exist(a array.Array) (ok bool) {
	a = a.Inverted()
	ok = it.Trie.Exist(a)
	return
}

// Match 判断 trie 是否匹配逆序的 a
func (it *InvertedTrie) Match(a array.Array) (ok bool) {
	a = a.Inverted()
	ok = it.Trie.Match(a)
	return
}
