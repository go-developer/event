// Package test...
//
// Author: go_developer@163.com<张德满>
//
// File:  etcd_driver_test.go.go
//
// Description: etcd_driver_test.go etcd 驱动的测试
//
// Date: 2020/10/10 11:58 上午
package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-developer/event/abstract"
	"github.com/go-developer/event/define"
	"github.com/go-developer/event/driver"
)

func getEtcdDriver() abstract.IDriver {
	ed, _ := driver.NewEtcdDriver(&define.EtcdDriverConfig{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5,
		Topic:       "event-queue",
		Buffer:      1,
	})
	return ed
}

// TestSubscribe ...
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:00 下午 2020/10/10
func TestSubscribe(t *testing.T) {
	ed := getRedisDriverInstance()
	go func() {
		for mes := range ed.Subscribe() {
			fmt.Println("订阅到的消息 : ", mes)
		}
	}()
	for i := 0; i < 10; i++ {
		ed.Publish(&define.Message{
			Type:      "etcd",
			Topic:     "event-queue",
			Key:       "event-queue",
			Driver:    "etcd",
			Timestamp: time.Now().Unix(),
			Data:      map[string]interface{}{"data": "etcd数据", "index": i},
		})
	}
}
