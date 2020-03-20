package setting

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BindKind(t *testing.T) {
	base := newSetting()
	s := newkindBindSetting(base)
	s.BindKind("testStr", reflect.String)

	r := s.checkKind("testStr", "testValue")
	assert.True(t, r)

	r = s.checkKind("testStr", 1)
	assert.False(t, r)
}

func Test_kindbindsetting_JSON(t *testing.T) {
	base0 := newSetting()
	m0 := newkindBindSetting(base0)

	m0.Set("123", 123)
	m0.Set("3.2.1", 321)

	base1 := newSetting()
	m1 := newkindBindSetting(base1)

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
