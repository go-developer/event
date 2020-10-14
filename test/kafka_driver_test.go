// Package test...
//
// Author: go_developer@163.com<张德满>
//
// File:  etcd_driver_test.go
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

func getKafkaDriver() abstract.IDriver {
	ed, err := driver.NewKafkaDriver(&define.KafkaDriverConfig{
		AddrList:        []string{"127.0.0.1:9092"},
		GroupID:         "local-test",
		PartitionCount:  3,
		Initial:         -2,
		Notifications:   false,
		CommitInterval:  100,
		ClientType:      define.KafkaClientTypeBoth,
		Topic:           "local-test",
		ProducerTimeout: 100,
		Buffer:          1024,
	})
	if nil != err {
		panic(err.Error())
	}
	return ed
}

// TestSubscribe ...
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:00 下午 2020/10/10
func TestKafkaSubscribe(t *testing.T) {
	ed := getKafkaDriver()
	go func() {
		for mes := range ed.Subscribe() {
			fmt.Println("订阅到的消息 : ", mes)
		}
	}()
	time.Sleep(5)
	for i := 0; i < 10; i++ {
		ed.Publish(&define.Message{
			Type:      "kafka",
			Topic:     "event-queue",
			Key:       fmt.Sprintf("event-queue-%d", i),
			Driver:    "etcd",
			Timestamp: time.Now().Unix(),
			Data:      map[string]interface{}{"data": "kafka数据", "index": i},
		})
	}
	time.Sleep(5 * time.Second)
}
