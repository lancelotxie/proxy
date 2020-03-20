package logger

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/lancelot/proxy/proxy.lib/ctxtransid"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	l := &logger{val: NOLOG}
	assert.NotNil(t, l)
	assert.Equal(t, int64(NOLOG), l.val)
	l.set(INFO)
	assert.Equal(t, int64(INFO), l.val)
	l.set(WARN)
	assert.Equal(t, int64(WARN), l.val)
	l.set(ERROR)
	assert.Equal(t, int64(ERROR), l.val)
	l.set(NOLOG)
	assert.Equal(t, int64(NOLOG), l.val)
	l.set(_LogLevel)
	assert.Equal(t, int64(_LogLevel), l.val)
}

func TestTrans2string(t *testing.T) {
	ctx := context.Background()
	actual := trans2string(ctx, infoheader, 1, 2, 3)
	assert.Equal(t, "123", actual)

	actual = trans2string(ctx, infoheader, 1, "2", []byte("345"))
	assert.Equal(t, "12[51 52 53]", actual)

	actual = trans2string(ctx, infoheader, "1", "\n", int64(3))
	assert.Equal(t, "1\n3", actual)

	actual = trans2string(ctx, infoheader, 1.2345, "6", 'N')
	assert.Equal(t, "1.2345678", actual)
}

type mockPrint struct {
	num     int32
	content []string
}

func newmockPrint(count int) *mockPrint {
	return &mockPrint{
		num:     0,
		content: make([]string, count),
	}
}
func (m *mockPrint) Println(opts ...interface{}) {
	m.content[m.num] = opts[0].(string)
	atomic.AddInt32(&m.num, 1)
}

func (m *mockPrint) Write(data []byte) (n int, err error) {
	m.content[m.num] = string(data)
	atomic.AddInt32(&m.num, 1)
	n = len(data)
	return
}

func (m *mockPrint) String() string {
	return fmt.Sprintf("%v\n%v\n", m.num, m.content)
}

func TestInfo(t *testing.T) {
	m := newmockPrint(100)
	l := &logger{val: 4}
	l.SetOutput(m)
	number := 0
	l.Set(INFO)
	ctx := ctxtransid.NewCTX(context.Background())
	l.Info(ctx, 1, 2, "3")
	assert.True(t, strings.Contains(m.content[number], infoheader))
	number++
	ctx2 := ctxtransid.NewCTX(context.Background())
	l.Info(ctx2, 1.2345, "6", 'N')
	assert.True(t, strings.Contains(m.content[number], infoheader))
}

func TestJudge(t *testing.T) {
	m := newmockPrint(100)
	l := &logger{val: 4}
	number := 0
	l.Set(ERROR)

	l.SetOutput(m)
	ctx := ctxtransid.NewCTX(context.Background())
	l.Info(ctx, 1, 2, "3")
	assert.False(t, strings.Contains(m.content[number], infoheader))
	assert.False(t, strings.Contains(m.content[number], errheader))

	ctx2 := ctxtransid.NewCTX(context.Background())
	l.Error(ctx2, 1.2345, "6", 'N')
	assert.False(t, strings.Contains(m.content[number], infoheader))
	assert.True(t, strings.Contains(m.content[number], errheader))

	ctx3 := ctxtransid.NewCTX(context.Background())
	l.Warning(ctx3, "1", 2, 'N')
	fmt.Println(m)
}

func TestInterface(t *testing.T) {
	ctx := ctxtransid.NewCTX(context.Background())
	Info(ctx, 1, 2, "3")
	ctx2 := ctxtransid.NewCTX(ctx)
	Error(ctx2, 1, 2345, "6", 'N')
	ctx3 := ctxtransid.NewCTX(ctx2)
	Warning(ctx3, "a", []byte("122345"), []int{1, 2, 3, 4, 5, 6})
}
