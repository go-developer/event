// Package define...
//
// Author: go_developer@163.com<张德满>
//
// File:  default.go
//
// Description: default 定义系统中各种默认值
//
// Date: 2020/9/23 8:57 下午
package define

const (
	// RedisDriverDefaultBuffer redis驱动消息缓冲区默认大小
	RedisDriverDefaultBuffer = 1
)

const (
	// KafkaDriverDefaultBuffer kafka驱动默认buffer缓冲区大小
	KafkaDriverDefaultBuffer = 1
	// KafkaDriverDefaultConsumerCount kafka驱动默认消费者数量
	KafkaDriverDefaultConsumerCount = 3
)

const (
	// EtcdDriverDefaultBuffer etcd驱动默认buffer大小
	EtcdDriverDefaultBuffer = 1
	// EtcdDriverDefaultTimeout etcd 驱动默认超时时间
	EtcdDriverDefaultTimeout = 1000
)
