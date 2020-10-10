// Package driver...
//
// Author: go_developer@163.com<张德满>
//
// File:  etcd.go
//
// Description: etcd 驱动的事件监听
//
// Date: 2020/10/9 6:42 下午
package driver

import (
	"context"
	"encoding/json"

	"github.com/coreos/etcd/clientv3"
	"github.com/go-developer/event/abstract"
	"github.com/go-developer/event/define"
)

// NewEtcdDriver etcd驱动实例
//
// Author : go_developer@163.com<张德满>
//
// Date : 11:29 上午 2020/10/10
func NewEtcdDriver(edc *define.EtcdDriverConfig) (abstract.IDriver, error) {
	ed := &etcdDriver{
		edc:           edc,
		messageChan:   make(chan *define.Message, edc.Buffer),
		exceptionChan: make(chan *define.Exception, edc.Buffer),
	}
	return ed, nil
}

type etcdDriver struct {
	edc           *define.EtcdDriverConfig // 配置文件
	client        *clientv3.Client         // etcd 客户端
	messageChan   chan *define.Message     // 消息管道
	exceptionChan chan *define.Exception   // 异常信息
}

// Init 初始化
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:49 下午 2020/10/9
func (ed *etcdDriver) Init() error {
	var (
		err error
	)
	if ed.client, err = clientv3.New(clientv3.Config{
		Endpoints:   ed.edc.Endpoints,
		DialTimeout: ed.edc.DialTimeout,
	}); nil != err {
		return err
	}
	return nil
}

// Publish 发布消息
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:50 下午 2020/10/9
func (ed *etcdDriver) Publish(mes *define.Message) error {
	var (
		kv clientv3.KV
		// resp     *clientv3.PutResponse
		err      error
		byteData []byte
	)
	if byteData, err = json.Marshal(mes); nil != err {
		return err
	}
	kv = clientv3.NewKV(ed.client)
	if _, err = kv.Put(context.TODO(), ed.edc.Topic, string(byteData)); nil != err {
		return err
	}
	return nil
}

// Subscribe 订阅消息
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:51 下午 2020/10/9
func (ed *etcdDriver) Subscribe() <-chan *define.Message {
	ed.startWatcher()
	return ed.messageChan
}

// SubscribeWithHandler 订阅消息并处理
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:51 下午 2020/10/9
func (ed *etcdDriver) SubscribeWithHandler(handler abstract.IHandler) {
	ed.startWatcher()
	for mes := range ed.messageChan {
		if err := handler.Handler(mes); nil != err {
			continue
		}
	}
}

// GetException 监听异常信息
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:52 下午 2020/10/9
func (ed *etcdDriver) GetException() <-chan *define.Exception {
	return nil
}

// startWatcher 启动事件行为监听
//
// Author : go_developer@163.com<张德满>
//
// Date : 11:45 上午 2020/10/10
func (ed *etcdDriver) startWatcher() {
	watchChan := ed.client.Watch(context.TODO(), ed.edc.Topic)
	go func() {
		for res := range watchChan {
			value := res.Events[0].Kv.Value
			var (
				message define.Message
				err     error
			)
			if err = json.Unmarshal(value, &message); nil != err {
				continue
			}
			ed.messageChan <- &message
		}
	}()
}
