package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetGet(t *testing.T) {
	c := newConfiguration()
	err := c.Set("1.2.3", "123")
	if !assert.NoError(t, err) {
		return
	}

	v, ok := c.GetString("1.2.3")
	if !assert.True(t, ok) {
		return
	}
	assert.Equal(t, "123", v)
}
