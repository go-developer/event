// Package test...
//
// Author: go_developer<张德满>
//
// File:  test_redis_driver.go
//
// Description: test_redis_driver
//
// Date: 2020/9/24 5:39 下午
package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-developer/event/define"
	"github.com/go-developer/event/driver"
)

// TestRedisDriver redis 驱动测试
//
// Author : go_developer@163.com<张德满>
//
// Date : 5:44 下午 2020/9/24
func TestRedisDriver(t *testing.T) {
	rdc := &define.RedisDriverConfig{
		Host:     "127.0.0.1",
		Port:     6379,
		Password: "",
		DB:       0,
		Timeout: struct {
			Connect int `json:"connect"`
			Write   int `json:"write"`
			Read    int `json:"read"`
		}{10, 10, 10},
		Buffer: 0,
	}
	rd := driver.NewRedisDriver(rdc)
	topic := "test-message"
	go func() {
		for mes := range rd.Subscribe(topic) {
			fmt.Println("订阅到的消息 :", mes)
		}
	}()
	for i := 0; i < 100; i++ {
		rd.Publish(&define.Message{
			Type:      "redis",
			Topic:     topic,
			Key:       topic,
			Driver:    "redis",
			Timestamp: time.Now().Unix(),
			Data:      map[string]interface{}{"message": fmt.Sprintf("测试数据 - %d", i)},
		})
		time.Sleep(1 * time.Second)
	}
}
