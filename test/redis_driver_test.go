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
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/go-developer/event/abstract"
	"github.com/go-developer/event/define"
	"github.com/go-developer/event/driver"
)

func getRedisDriverInstance() abstract.IDriver {
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
	return rd
}

func sendMessage(rd abstract.IDriver, topic string, cnt int) {
	for i := 0; i < cnt; i++ {
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

// TestRedisDriver redis 驱动测试
//
// Author : go_developer@163.com<张德满>
//
// Date : 5:44 下午 2020/9/24
func TestRedisDriver(t *testing.T) {
	rd := getRedisDriverInstance()
	topic := "test-message"
	go func() {
		for mes := range rd.Subscribe(topic) {
			fmt.Println("订阅到的消息 :", mes)
		}
	}()
	sendMessage(rd, topic, 10)
}

// TestRedisDriverWithHandler 测试带handler的消息订阅
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:19 下午 2020/10/9
func TestRedisDriverWithHandler(t *testing.T) {
	rd := getRedisDriverInstance()
	topic := "test-message"
	handler := new(testRedisDriverWithHandler)
	rd.StartSubscribeWithHandler(topic, handler)
	sendMessage(rd, topic, 10)
	go func() {
		for {
			select {
			case err := <-rd.GetException():
				fmt.Println("处理异常 : " + err.Message)
			}
		}
	}()
}

type testRedisDriverWithHandler struct {
}

func (trdh *testRedisDriverWithHandler) Handler(message *define.Message) error {
	if byteData, err := json.Marshal(message); nil != err {
		return err
	} else {
		fmt.Println("通过handler输出的信息 : ", string(byteData))
		return nil
	}
}
