// Package define...
//
// Author: go_developer@163.com<张德满>
//
// File:  driver.go
//
// Description: driver 定义支持的事件驱动
//
// Date: 2020/9/23 3:24 下午
package define

import (
	"time"
)

const (
	// EventDriverRedis redis驱动
	EventDriverRedis = "redis"
	// EventDriverApollo apollo驱动
	EventDriverApollo = "apollo"
	// EventDriverZookeeper zk驱动
	EventDriverZookeeper = "zookeeper"
	// EventDriverEtcd etcd驱动
	EventDriverEtcd = "etcd"
	// EventDriverKafka kafka事件驱动
	EventDriverKafka = "kafka"
)

// DriverConfig 定于驱动的配置
//
// Author : go_developer@163.com<张德满>
//
// Date : 3:31 下午 2020/9/23
type DriverConfig struct {
	Type     string `json:"type"`     // 驱动类型
	Host     string `json:"host"`     // 主机
	Port     string `json:"port"`     // 端口
	Username string `json:"username"` // 账号
	Password string `json:"password"` // 密码
}

// RedisDriverConfig redis驱动的配置
//
// Author : go_developer@163.com<张德满>
//
// Date : 8:49 下午 2020/9/23
type RedisDriverConfig struct {
	Host     string `json:"host"`     // 主机
	Port     int    `json:"port"`     // 端口
	Password string `json:"password"` // 密码
	DB       int    `json:"db"`       // db
	Topic    string `json:"topic"`    // 订阅的消息队列
	Timeout  struct {
		Connect int `json:"connect"`
		Write   int `json:"write"`
		Read    int `json:"read"`
	} `json:"timeout"` // 超时相关配置
	Buffer int64 `json:"buffer"` // 消息订阅缓冲区大小
}

// KafkaDriverConfig kafka驱动的配置
//
// Author : go_developer@163.com<张德满>
//
// Date : 2:15 下午 2020/10/9
type KafkaDriverConfig struct {
	AddrList        []string      `json:"addr_list"`        // 集群地址
	GroupID         string        `json:"group_id"`         // 消费者组ID
	PartitionCount  int           `json:"partition_count"`  // 分区数量
	Initial         int64         `json:"initial"`          // 从最新偏移量开始消费 or 最老 -1 最新 -2 最老
	Notifications   bool          `json:"notifications"`    // 返回通知信息
	CommitInterval  time.Duration `json:"commit_interval"`  // 自动提交offset时间间隔
	ClientType      int           `json:"client_type"`      // 类型 0 - 仅生产 1 - 仅消费 2 - 生产 & 消费
	Topic           string        `json:"topic"`            // 订阅的topic
	ProducerTimeout time.Duration `json:"producer_timeout"` // 生产消息超时时间
	Buffer          int64         `json:"buffer"`           // 消息订阅缓冲区大小
}

// 定义kafka client的类型
//
// Author : go_developer@163.com<张德满>
//
// Date : 3:04 下午 2020/10/9
const (
	// KafkaClientTypeProduce 生产者客户端
	KafkaClientTypeProduce = 1
	// KafkaClientTypeConsumer 消费者客户端
	KafkaClientTypeConsumer = 2
	// 同时作为 生产者 & 消费者客户端
	KafkaClientTypeBoth = 3
)

// EtcdDriverConfig 基于etcd的事件发布订阅
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:32 下午 2020/10/9
type EtcdDriverConfig struct {
	Endpoints   []string      `json:"endpoints"`    // etcd节点信息
	DialTimeout time.Duration `json:"dial_timeout"` // 超时配置
	Topic       string        `json:"topic"`        // 监听的topic
	Buffer      int           `json:"buffer"`       // 数据缓冲区大小
}
