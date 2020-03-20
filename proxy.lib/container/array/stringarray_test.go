package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringArray(t *testing.T) {
	var sa Array = StringArray("test")

	e := sa.Get(1)
	if !assert.Equal(t, 'e', int32(e.(byte))) {
		return
	}

	_t := sa.Get(3)
	if !assert.Equal(t, 't', int32(_t.(byte))) {
		return
	}

	l := sa.Len()
	assert.Equal(t, 4, l)
}
