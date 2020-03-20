package rpcconfig

import "encoding/json"

// KeyValue 键值对
type KeyValue struct {
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// KeyValues 键值对集合
type KeyValues []KeyValue

// Push 添加键值对
func (kvs KeyValues) Push(kv KeyValue) (r KeyValues) {
	r = append(kvs, kv)
	return
}

// Bytes 转换为 bytes
func (kvs KeyValues) Bytes() []byte {
	j, _ := json.Marshal(kvs)
	return j
}

// ParseKeyValues 解析键值对集合
func ParseKeyValues(b []byte) (kvs KeyValues, err error) {
	err = json.Unmarshal(b, &kvs)
	return
}
