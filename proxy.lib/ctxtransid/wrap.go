package ctxtransid

import (
	"context"
	"strconv"
)

// NewCTX : 根据context，返回带有特定key-value 的子context,该key不可指定
// 每NewCTX一次，该key的值都会发生变化
func NewCTX(CTX context.Context) (newCTX context.Context) {
	transactionID := int(transactionIDCreator.Up())
	newCTX = context.WithValue(CTX, keyTransactionID, transactionID)
	return
}

// GetID : 根据context，返回 特定key的value, 该key不可指定
func GetID(CTX context.Context) (id string) {
	// 让ctx为nil时，程序panic
	// if ctx == nil {
	// return
	// }
	if v := CTX.Value(keyTransactionID); v != nil {
		id = strconv.Itoa(v.(int))
	}
	return
}

// WithID : 根据指定 值，返回特定key 的context。该key不可指定
func WithID(CTX context.Context, id int) (newCTX context.Context) {
	newCTX = context.WithValue(CTX, keyTransactionID, id)
	return
}
