package base

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Resolver(t *testing.T) {
	resolv := func(in string) (out interface{}, err error) {
		out = in
		return
	}

	r := New(resolv)
	out, _ := r.Resolv("1")
	assert.Equal(t, "1", out)
}
