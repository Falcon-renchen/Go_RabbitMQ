package main

import (
	"Go_RabbitMQ/rabbitmq03/mq"
	"fmt"
)

type TestPro struct {
	msgContent string
}

//实现发送者
func (t *TestPro) MsgContent() string {
	return t.msgContent
}

//实现接收者
func (t *TestPro) Consumer(dataByte []byte) error {
	fmt.Println(string(dataByte))
	return nil
}

func main() {
	msg := fmt.Sprintf("这是测试任务")
	t := &TestPro{
		msgContent: msg,
	}

	queueExchange := &mq.QueueExchange{
		"test.rabbit",
		"rabbit.key",
		"test.rabbit.mq",
		"direct",
	}
	mq := mq.New(queueExchange)
	mq.RegisterProducer(t)
	mq.RegisterReceiver(t)
	mq.Start()
}
