// Package abstract...
//
// Author: go_developer@163.com<张德满>
//
// File:  driver.go
//
// Description: driver 驱动的接口定义
//
// Date: 2020/9/23 3:39 下午
package abstract

import (
	"github.com/go-developer/event/define"
)

// IDriver 定义驱动的接口
//
// Author : zhangdeman001@ke.com<张德满>
//
// Date : 3:43 下午 2020/9/23
type IDriver interface {
	// Init 初始化驱动
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 3:53 下午 2020/9/23
	Init() error
	// Publish 发布事件消息
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 7:51 下午 2020/9/23
	Publish(message *define.Message) error

	// Subscribe 订阅事件消息
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 7:51 下午 2020/9/23
	Subscribe(topic string) <-chan *define.Message

	// GetException 获取异常信息
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 12:06 下午 2020/9/24
	GetException() <-chan *define.Exception
}
