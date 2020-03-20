package ctxtransid

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	CTX := context.Background()
	newc := NewCTX(CTX)
	id := GetID(newc)
	assert.NotEmpty(t, id)

	v := CTX.Value(keyTransactionID)
	fmt.Println(v == nil)
	id2 := GetID(CTX)
	assert.Empty(t, id2)

	c3 := NewCTX(newc)
	id3 := GetID(c3)
	expected := strconv.Itoa(int(transactionIDCreator.Value))
	assert.Equal(t, expected, id3)
}

func TestGetID(t *testing.T) {
	id2 := GetID(context.WithValue(context.Background(), context.Background(), 4))
	assert.Empty(t, id2)
	id3 := GetID(NewCTX(context.Background()))
	assert.NotEmpty(t, id3)
}
