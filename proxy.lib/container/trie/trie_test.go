package trie

import (
	"testing"

	"github.com/lancelot/proxy/proxy.lib/container/array"

	"github.com/stretchr/testify/assert"
)

func Test_SetExist(t *testing.T) {
	_t := New()

	a0 := array.StringArray("test0")
	a1 := array.StringArray("test1")
	a2 := array.StringArray("test")
	a3 := array.StringArray("test01")

	_t.Set(a0)

	if !assert.True(t, _t.Exist(a0)) {
		return
	}
	if !assert.False(t, _t.Exist(a1)) {
		return
	}
	if !assert.False(t, _t.Exist(a2)) {
		return
	}
	if !assert.False(t, _t.Exist(a3)) {
		return
	}
}

func Test_Match(t *testing.T) {
	_t := New()

	a0 := array.StringArray("test0")
	a1 := array.StringArray("test1")
	a2 := array.StringArray("test")
	a3 := array.StringArray("test01")

	_t.Set(a0)

	if !assert.True(t, _t.Match(a0)) {
		return
	}
	if !assert.False(t, _t.Match(a1)) {
		return
	}
	if !assert.False(t, _t.Match(a2)) {
		return
	}
	if !assert.True(t, _t.Match(a3)) {
		return
	}
}
