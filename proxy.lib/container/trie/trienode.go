package trie

type trieNode struct {
	colored bool // 该节点是否上色，表示是否存在以该节点结束的串
	childs  map[interface{}]*trieNode
}

func newNode() (n *trieNode) {
	n = new(trieNode)
	n.childs = make(map[interface{}]*trieNode)
	return
}

func (tn *trieNode) find(key interface{}) (child *trieNode, ok bool) {
	child, ok = tn.childs[key]
	return
}

func (tn *trieNode) grow(key interface{}) (child *trieNode) {
	child = newNode()
	tn.childs[key] = child
	return
}

// forceGetChild 强制获取子节点，如果获取失败则创建
func (tn *trieNode) forceGetChild(key interface{}) (child *trieNode) {
	child, ok := tn.find(key)
	if !ok {
		child = tn.grow(key)
	}
	return
}
