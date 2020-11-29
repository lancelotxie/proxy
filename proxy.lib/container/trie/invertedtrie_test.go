package trie

import (
	"testing"

	"github.com/lancelotXie/proxy/proxy.lib/container/array"

	"github.com/stretchr/testify/assert"
)

func Test_InvertedSetExist(t *testing.T) {
	_t := New()

	a0 := array.StringArray("test0")
	a1 := array.StringArray("test1")

	_t.Set(a0)

	if !assert.True(t, _t.Exist(a0)) {
		return
	}

	assert.False(t, _t.Exist(a1))
}
