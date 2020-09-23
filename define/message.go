// Package define...
//
// Author: go_developer@163.com<张德满>
//
// File:  message.go
//
// Description: message 消息结构体的定义
//
// Date: 2020/9/23 7:53 下午
package define

// Message 定义消息的数据结构
//
// Author : go_developer@163.com<张德满>
//
// Date : 7:53 下午 2020/9/23
type Message struct {
	Type      string                 `json:"type"`      // 消息类型
	Topic     string                 `json:"topic"`     // 订阅的主题
	Key       string                 `json:"key"`       // 用于数据分区(仅针对kafka有效)
	Driver    string                 `json:"driver"`    // 驱动类型
	Timestamp int64                  `json:"timestamp"` // 时间戳
	Data      map[string]interface{} `json:"data"`      // 数据
}
