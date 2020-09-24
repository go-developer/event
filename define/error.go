// Package define...
//
// Author: go_developer@163.com<张德满>
//
// File:  error.go
//
// Description: error 定义异常信息
//
// Date: 2020/9/24 11:45 上午
package define

const (
	// ExceptionTypeFormatError 消息数据格式错误
	ExceptionTypeFormatError = "data format error"
	ExceptionType
)

// Exception 自定义异常
//
// Author : go_developer@163.com<张德满>
//
// Date : 11:50 上午 2020/9/24
type Exception struct {
	Type          string `json:"type"`           // 异常类型
	Code          int    `json:"code"`           // 异常code
	SubscribeData string `json:"subscribe_data"` // 订阅到的原始数据
}
