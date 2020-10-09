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
	"github.com/go-developer/event/abstract"
	"github.com/go-developer/event/define"
)

func NewEtcdDriver(edc *define.EtcdDriverConfig) (abstract.IDriver, error) {
	ed := &etcdDriver{edc: edc}
	return ed, nil
}

type etcdDriver struct {
	edc *define.EtcdDriverConfig // 配置文件
}

// Init 初始化
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:49 下午 2020/10/9
func (ed *etcdDriver) Init() error {
	return nil
}

// Publish 发布消息
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:50 下午 2020/10/9
func (ed *etcdDriver) Publish(mes *define.Message) error {
	return nil
}

// Subscribe 订阅消息
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:51 下午 2020/10/9
func (ed *etcdDriver) Subscribe() <-chan *define.Message {
	return nil
}

// SubscribeWithHandler 订阅消息并处理
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:51 下午 2020/10/9
func (ed *etcdDriver) SubscribeWithHandler(handler abstract.IHandler) {
	return
}

// GetException 监听异常信息
//
// Author : go_developer@163.com<张德满>
//
// Date : 6:52 下午 2020/10/9
func (ed *etcdDriver) GetException() <-chan *define.Exception {
	return nil
}
