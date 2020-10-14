// Package driver...
//
// Author: go_developer@163.com<张德满>
//
// File:  kafka.go
//
// Description: kafka 基于kafka实现的消息订阅与发布
//
// Date: 2020/10/9 12:36 下午
package driver

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-developer/event/abstract"
	"github.com/go-developer/event/define"
	"github.com/go-developer/go-util/util"
)

// NewKafkaDriver 获取kafka的驱动
//
// Author : go_developer@163.com<张德满>
//
// Date : 2:29 下午 2020/10/9
func NewKafkaDriver(kdc *define.KafkaDriverConfig) (abstract.IDriver, error) {
	kd := &kafkaDriver{kdc: kdc}
	if err := kd.Init(); nil != err {
		return nil, err
	}
	return kd, nil
}

type kafkaDriver struct {
	kdc           *define.KafkaDriverConfig // 配置
	consumer      []*cluster.Consumer       // 消费者实例
	producer      sarama.SyncProducer       // 生产者实例
	messageChan   chan *define.Message      // 消息管道
	exceptionChan chan *define.Exception    // 异常信息
}

// Init kafka初始化
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:38 下午 2020/10/9
func (kd *kafkaDriver) Init() error {
	// 既不是生产者，也不是消费者
	if !kd.isProducerClient() && !kd.isConsumerClient() {
		return fmt.Errorf(
			"请指定正确的 client type %d - 生产者客户端 %d - 消费者客户端 %d - 即为生产者客户端，也是消费者客户端",
			define.KafkaClientTypeProduce, define.KafkaClientTypeConsumer, define.KafkaClientTypeBoth,
		)
	}
	var (
		err error
	)
	// 初始化生产者
	if kd.isProducerClient() {
		if err = kd.initProducer(); nil != err {
			return err
		}
	}
	// 初始化消费者
	if kd.isConsumerClient() {
		if err = kd.initConsumer(); nil != err {
			return err
		}
	}
	return nil
}

// Publish 发布kafka消息
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:39 下午 2020/10/9
func (kd *kafkaDriver) Publish(message *define.Message) error {
	byteData, err := json.Marshal(message)
	if nil != err {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: kd.kdc.Topic,
		Value: sarama.ByteEncoder(string(byteData)),
	}
	if _, _, err := kd.producer.SendMessage(msg); nil != err {
		return err
	}
	return nil
}

// Subscribe 订阅kafka消息 ...
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:39 下午 2020/10/9
func (kd *kafkaDriver) Subscribe() <-chan *define.Message {
	// 启动消费者
	kd.startConsumer()
	return kd.messageChan
}

// SubscribeWithHandler 带handler的消息订阅
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:41 下午 2020/10/9
func (kd *kafkaDriver) SubscribeWithHandler(handler abstract.IHandler) {
	for msg := range kd.Subscribe() {
		if err := handler.Handler(msg); nil != err {
			byteData, _ := json.Marshal(msg)
			kd.setExceptionInfo(define.ExceptionTypeHandlerError, err.Error(), string(byteData))
		}
	}
}

// GetException 异常获取
//
// Author : go_developer@163.com<张德满>
//
// Date : 12:41 下午 2020/10/9
func (kd *kafkaDriver) GetException() <-chan *define.Exception {
	return kd.exceptionChan
}

// isProducerClient 判断是否为生产者客户端
//
// Author : go_developer@163.com<张德满>
//
// Date : 3:07 下午 2020/10/9
func (kd *kafkaDriver) isProducerClient() bool {
	return kd.kdc.ClientType == define.KafkaClientTypeProduce || kd.kdc.ClientType == define.KafkaClientTypeBoth
}

// isConsumerClient 判断是否为消费者客户端
//
// Author : go_developer@163.com<张德满>
//
// Date : 3:07 下午 2020/10/9
func (kd *kafkaDriver) isConsumerClient() bool {
	return kd.kdc.ClientType == define.KafkaClientTypeConsumer || kd.kdc.ClientType == define.KafkaClientTypeBoth
}

// initConsumer 初始化消费者客户端
//
// Author : go_developer@163.com<张德满>
//
// Date : 3:20 下午 2020/10/9
func (kd *kafkaDriver) initConsumer() error {
	config := cluster.NewConfig()
	config.Group.Return.Notifications = kd.kdc.Notifications
	config.Consumer.Offsets.CommitInterval = kd.kdc.CommitInterval
	config.Consumer.Offsets.Initial = kd.kdc.Initial
	if kd.kdc.PartitionCount <= 0 {
		kd.kdc.PartitionCount = define.KafkaDriverDefaultConsumerCount
	}
	if kd.kdc.Buffer <= 0 {
		kd.kdc.Buffer = define.KafkaDriverDefaultBuffer
	}
	kd.consumer = make([]*cluster.Consumer, 0)
	// 生成GroupID,确保每台机器的GroupID不一样,每台机器都可以订阅到消息
	serverIP, err := util.ProjectUtil.GetServerIP()
	if nil != err {
		return err
	}
	groupID := fmt.Sprintf("%s - %s", serverIP, kd.kdc.GroupID)
	// 多少个分区就启动多少个消费者
	for i := 0; i < kd.kdc.PartitionCount; i++ {
		c, err := cluster.NewConsumer(
			kd.kdc.AddrList,
			groupID,
			[]string{kd.kdc.Topic},
			config,
		)
		if nil != err {
			return err
		}
		kd.consumer = append(kd.consumer, c)
	}
	kd.messageChan = make(chan *define.Message, kd.kdc.Buffer)
	return nil
}

// initProducer 初始化生产者
//
// Author : go_developer@163.com<张德满>
//
// Date : 3:23 下午 2020/10/9
func (kd *kafkaDriver) initProducer() error {
	var (
		err error
	)
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Timeout = kd.kdc.ProducerTimeout
	if kd.producer, err = sarama.NewSyncProducer(kd.kdc.AddrList, config); nil != err {
		return err
	}
	return nil
}

// startConsumer 启动消费者
//
// Author : go_developer@163.com<张德满>
//
// Date : 4:36 下午 2020/10/9
func (kd *kafkaDriver) startConsumer() {
	for _, c := range kd.consumer {
		go func(consumer *cluster.Consumer) {
			for msg := range consumer.Messages() {
				offset := strconv.FormatInt(msg.Offset, 10)
				consumer.MarkOffset(msg, offset)
				var data define.Message
				if err := json.Unmarshal(msg.Value, &data); nil != err {
					kd.setExceptionInfo(define.ExceptionTypeFormatError, "解析订阅道德kafka消息失败", string(msg.Value))
					continue
				}
				kd.messageChan <- &data
			}
		}(c)
	}
}

// setExceptionInfo 设置异常信息
//
// Author : go_developer@163.com<张德满>
//
// Date : 4:32 下午 2020/10/9
func (kd *kafkaDriver) setExceptionInfo(exceptionType string, message string, data string) {
	go func() {
		kd.exceptionChan <- &define.Exception{
			Type:          exceptionType,
			Message:       message,
			SubscribeData: data,
		}
	}()
}
