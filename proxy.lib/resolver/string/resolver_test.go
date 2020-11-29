package string

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lancelotXie/proxy/proxy.lib/resolver/base"
)

func Test_Resolv(t *testing.T) {
	resolv := func(in string) (out interface{}, err error) {
		if in == "2" {
			out = 2
		} else {
			out = in
		}
		return
	}
	br := base.New(resolv)
	ss := New(br)

	out, err := ss.Resolv("1")
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, "1", out) {
		return
	}

	_, err = ss.Resolv("2")
	assert.Error(t, err)
}
