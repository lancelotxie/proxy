package ctxtransid

import (
	"github.com/eaglexiang/go/counter"
)

var transactionIDCreator = counter.Counter{Value: 0}

// contextKey :
type contextKey int

// keyTransactionID :
const (
	keyTransactionID contextKey = iota
)
