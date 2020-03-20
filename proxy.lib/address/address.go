package address

import (
	"net"
	"strconv"
	"strings"
)

// Unpack : 仅支持ipv4,
// bug ：ip 为非法ip时，依然返回ok：true
func Unpack(addr net.Addr) (ip string, port int, network string, ok bool) {
	if addr == nil {
		return "", 0, "", false
	}
	ipe := addr.String()
	n := strings.LastIndex(ipe, ":")
	if n < 1 {
		return "", 0, "", false
	}
	ip = ipe[:n]
	n++ // ":" 会占用一个字符
	port, err := strconv.Atoi(ipe[n:])
	if err != nil {
		return "", 0, "", false
	}
	network = addr.Network()
	ok = true
	return
}

// Address : net.conn 所需要的net.Addr
type Address struct {
	ip      string
	port    int
	network string
}

// New : 返回 Address。已实现net.Addr 接口
func New(ip string, p int, net string) net.Addr {
	return Address{ip: ip, port: p, network: net}
}

// Network :
func (a Address) Network() string { return a.network }

// String :
func (a Address) String() string { return a.ip + ":" + strconv.Itoa(a.port) }
