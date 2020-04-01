package streamservice

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// ReqAddr : DualStream流传输中，用于传递ip，port，network，id的结构体
type ReqAddr struct {
	IP      string `json:"ip,omitempty"`
	Port    int    `json:"Port,omitempty"`
	Network string `json:"Network,omitempty"`
	ID      string `json:"ID,omitempty"` //context的键id的值
}

// NewReqAddr : 构造ReqAddr
func NewReqAddr(ip string, port int, network string, id string) *ReqAddr {
	return &ReqAddr{
		IP:      ip,
		Port:    port,
		Network: network,
		ID:      id,
	}
}

// NewStreamBytesWithoutData : 构造双向传输中 只传 目标地址的 包
func NewStreamBytesWithoutData(ip string, port int, network string, id string) (req *StreamBytes, err error) {
	connInfo := NewReqAddr(ip, port, network, id)
	info, err := json.Marshal(connInfo)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	req = &StreamBytes{
		Data: []byte(info),
	}
	return
}

// UnpackInfo : 解析 双向传输中 只传 地址 的包,解析出 地址中的 ip,port,network
func UnpackInfo(info StreamBytes) (ip string, port int, network string, id string, err error) {
	var addrInfo ReqAddr
	err = json.Unmarshal(info.Data, &addrInfo)
	if err != nil {
		err = errors.WithStack(err)
	}
	ip = addrInfo.IP
	port = addrInfo.Port
	network = addrInfo.Network
	id = addrInfo.ID
	return
}

// NewStreamBytes : 构造 双向传输 中只传数据的 包
func NewStreamBytes(data []byte) *StreamBytes {
	return &StreamBytes{
		Data: data,
	}
}

// NewDomainReq : 构造 解析ip地址的请求
func NewDomainReq(domain string) *DomainReq {
	return &DomainReq{Domain: domain}
}

// NewDomainResponse : 构造 解析ip地址的回复
func NewDomainResponse(ip string) *IPRespose {
	return &IPRespose{IP: ip}
}

// NewLocationReq : 构造 解析ip地址的请求
func NewLocationReq(ip string) *LocationReq {
	return &LocationReq{IP: ip}
}

// NewLocationResponse : 构造 解析ip地址的回复
func NewLocationResponse(location string) *LocationResponse {
	return &LocationResponse{Location: location}
}
