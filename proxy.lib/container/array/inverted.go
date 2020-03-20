package array

// InvertedArray 逆序的数组
type InvertedArray struct {
	Array
}

// Get 根据逆序的 index 获取 value
func (isa InvertedArray) Get(index int) (value interface{}) {
	// 根据逆序的 index 获取实际的 index
	lastIndex := isa.Len() - 1
	index = lastIndex - index

	value = isa.Array.Get(index)
	return
}

func newInvertedArray(a Array) (inverted Array) {
	_a := new(InvertedArray)
	_a.Array = a
	inverted = _a
	return
}
