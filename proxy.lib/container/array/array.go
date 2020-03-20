package array

// Array 数组
type Array interface {
	Get(index int) (value interface{})
	Len() (len int)
	Inverted() Array
}
