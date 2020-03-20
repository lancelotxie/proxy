package array

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Inverted(t *testing.T) {
	var a Array = StringArray("1234")
	var inverted = a.Inverted()

	if !assert.Equal(t, '1', int32(a.Get(0).(byte))) {
		return
	}

	assert.Equal(t, '4', int32(inverted.Get(0).(byte)))
}
