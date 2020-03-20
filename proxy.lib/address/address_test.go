package address

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpack(t *testing.T) {
	ip := "144.144.144.144"
	port := 8090
	network := "tcp"

	addr := New(ip, port, network)
	outip, outport, outnet, ok := Unpack(addr)
	assert.True(t, ok)
	assert.Equal(t, ip, outip)
	assert.Equal(t, outport, port)
	assert.Equal(t, outnet, network)
}
