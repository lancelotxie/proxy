package setting

import (
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Hook(t *testing.T) {
	s := newSetting()
	ks := newkindBindSetting(s)
	hs := newhookSetting(ks)

	var signal string
	h := Hooker{
		H0: func(value interface{}) (err error) {
			if value != "testValue" {
				err = errors.New("testError")
			}
			return
		},
		H1: func(value interface{}, ok bool) {
			if ok {
				signal = value.(string)
			}
		},
	}

	hs.Hook("testKey", h)

	err := hs.Set("testKey", "notTestValue")
	assert.Error(t, err)
	assert.NotEqual(t, "notTestValue", signal)

	err = hs.Set("testKey", "testValue")
	assert.NoError(t, err)
	assert.Equal(t, "testValue", signal)
}

func Test_hooksetting_JSON(t *testing.T) {
	s0 := newSetting()
	ks0 := newkindBindSetting(s0)
	m0 := newhookSetting(ks0)

	m0.Set("123", 123)
	m0.Set("3.2.1", 321)

	s1 := newSetting()
	ks1 := newkindBindSetting(s1)
	m1 := newhookSetting(ks1)

	b, err := json.Marshal(m0)
	assert.NoError(t, err)

	err = json.Unmarshal(b, &m1)
	assert.NoError(t, err)

	v0, ok := m1.Get("123")
	assert.True(t, ok)
	assert.Equal(t, float64(123), v0)

	v1, ok := m1.Get("3.2.1")
	assert.True(t, ok)

	assert.Equal(t, float64(321), v1)
}
