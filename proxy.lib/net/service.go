package net

// Service 可开启可关闭的服务
type Service interface {
	Start() error
	Close() error
}
