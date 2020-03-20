package array

// StringArray array 形式的 string
type StringArray string

// Get 获取 index 对应的 byte
func (sa StringArray) Get(index int) (value interface{}) {
	value = sa[index]
	return
}

// Len string 作为 []byte 的长度
func (sa StringArray) Len() (l int) {
	l = len(sa)
	return
}

// String 输出 string 原型
func (sa StringArray) String() string {
	return string(sa)
}

// Inverted 输出逆序版 array
func (sa StringArray) Inverted() (inverted Array) {
	inverted = newInvertedArray(sa)
	return
}
