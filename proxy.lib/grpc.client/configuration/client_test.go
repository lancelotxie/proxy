package configuration

import (
	"context"
	"testing"
	"time"

	"github.com/lancelot/proxy/proxy.lib/configuration/grpc"
	"github.com/lancelot/proxy/proxy.lib/path"

	"github.com/stretchr/testify/assert"
)

func Test_SetGet(t *testing.T) {
	path.Init(path.Client)

	go func() {
		err := grpc.Start("127.0.0.1", 8086)
		if !assert.NoError(t, err) {
			return
		}
	}()
	defer grpc.Stop()

	time.Sleep(time.Millisecond * 100)
	err := Init("127.0.0.1", 8086)
	if !assert.NoError(t, err) {
		return
	}

	ctx := context.Background()
	err = Set(ctx, "testKey", "testValue")
	if !assert.NoError(t, err) {
		return
	}

	v, err := Get(ctx, "testKey")
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "testValue", v)
}
