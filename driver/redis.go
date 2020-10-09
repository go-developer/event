// Package driver...
//
// Author: go_developer@163.com<张德满>
//
// File:  redis.go
//
// Description: redis 驱动的初始化
//
// Date: 2020/9/23 8:27 下午
package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-developer/event/abstract"
	"github.com/go-developer/event/define"
	"github.com/go-redis/redis/v8"
)

// NewRedisDriver 获取 redis driver实例
//
// Author : go_developer@163.com<张德满>
//
// Date : 8:32 下午 2020/9/23
func NewRedisDriver(cf *define.RedisDriverConfig) (abstract.IDriver, error) {
	rd := &redisDriver{
		cf: cf,
	}
	if err := rd.Init(); nil != err {
		return nil, err
	}
	return rd, nil
}

type redisDriver struct {
	cf            *define.RedisDriverConfig // redis 配置
	messageChan   chan *define.Message      // 数据缓冲chan
	instance      *redis.Client             // redis 缓冲chan
	exceptionChan chan *define.Exception    // 异常信息队列
}

// Init 初始化
//
// Author : zhangdeman001@ke.com<张德满>
//
// Date : 8:29 下午 2020/9/23
func (rd *redisDriver) Init() error {
	option := &redis.Options{
		DB:           rd.cf.DB,
		Addr:         fmt.Sprintf("%s:%d", rd.cf.Host, rd.cf.Port),
		DialTimeout:  time.Duration(rd.cf.Timeout.Connect*1000) * time.Millisecond,
		ReadTimeout:  time.Duration(rd.cf.Timeout.Read*1000) * time.Millisecond,
		WriteTimeout: time.Duration(rd.cf.Timeout.Read*1000) * time.Millisecond,
		Password:     rd.cf.Password,
	}
	rd.instance = redis.NewClient(option)
	if nil == rd.instance {
		panic("connect to redis server fail")
	}
	if rd.cf.Buffer <= 0 {
		rd.cf.Buffer = define.RedisDriverDefaultBuffer
	}
	rd.messageChan = make(chan *define.Message, rd.cf.Buffer)
	rd.exceptionChan = make(chan *define.Exception, rd.cf.Buffer)
	return nil
}

// Publish 发布事件
//
// Author : go_developer@163.com<张德满>
//
// Date : 8:29 下午 2020/9/23
func (rd *redisDriver) Publish(message *define.Message) error {
	byteData, _ := json.Marshal(message)
	if err := rd.instance.Publish(context.Background(), rd.cf.Topic, string(byteData)).Err(); nil != err {
		return err
	}
	return nil
}

// subscribe 订阅消息
//
// Author : go_developer@163.com<张德满>
//
// Date : 8:31 下午 2020/9/23
func (rd *redisDriver) Subscribe() <-chan *define.Message {
	go func() {
		// 拉取数据
		pubSubRes := rd.instance.Subscribe(context.Background(), rd.cf.Topic)
		for mes := range pubSubRes.Channel() {
			var (
				mesData define.Message
				err     error
			)
			if err = json.Unmarshal([]byte(mes.Payload), &mesData); nil != err {
				rd.setExceptionInfo(define.ExceptionTypeFormatError, err.Error(), mes.Payload)
				continue
			}
			rd.messageChan <- &mesData
		}
	}()

	return rd.messageChan
}

// StartSubscribeWithHandler 支持直接处理数据
//
// Author : go_developer@163.com<张德满>
//
// Date : 2:12 下午 2020/9/24
func (rd *redisDriver) SubscribeWithHandler(handler abstract.IHandler) {
	go func() {
		// 拉取数据
		pubSubRes := rd.instance.Subscribe(context.Background(), rd.cf.Topic)
		for mes := range pubSubRes.Channel() {
			var (
				mesData define.Message
				err     error
			)
			if err = json.Unmarshal([]byte(mes.Payload), &mesData); nil != err {
				rd.setExceptionInfo(define.ExceptionTypeFormatError, err.Error(), mes.Payload)
				continue
			}
			if err = handler.Handler(&mesData); nil != err {
				rd.setExceptionInfo(define.ExceptionTypeHandlerError, err.Error(), mes.Payload)
			}
		}
	}()
}

// GetException 异常信息
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:06 下午 2020/9/24
func (rd *redisDriver) GetException() <-chan *define.Exception {
	return rd.exceptionChan
}

// setExceptionInfo 设置异常信息
//
// Author : zhangdeman001@ke.com<张德满>
//
// Date : 2:15 下午 2020/9/24
func (rd *redisDriver) setExceptionInfo(exceptionType string, message string, data string) {
	go func() {
		rd.exceptionChan <- &define.Exception{
			Type:          exceptionType,
			Message:       message,
			SubscribeData: data,
		}
	}()
}
