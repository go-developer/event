// Package abstract...
//
// Author: <张德满>
//
// File:  handler.go
//
// Description: handler
//
// Date: 2020/9/24 12:55 下午
package abstract

import (
	"github.com/go-developer/event/define"
)

// IHandler 数据处理接口的定义
//
// Author : go_developer@163.com<张德满>
//
// Date : 2:07 下午 2020/9/24
type IHandler interface {
	// Handler 处理得到的消息
	//
	// Author : go_developer@163.com<张德满>
	//
	// Date : 2:08 下午 2020/9/24
	Handler(message *define.Message) error
}
