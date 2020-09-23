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
