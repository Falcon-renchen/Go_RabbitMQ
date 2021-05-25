package main

import (
	"Go_RabbitMQ/rabbitmq07_demo/Lib"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//假设这是真正的发送邮件的函数
func Send(c string, msg amqp.Delivery) error {
	time.Sleep(time.Second * 3) //假设很耗时,用协程
	fmt.Printf("%s向userID=%s的用户发送邮件\n", c, string(msg.Body))
	msg.Ack(false)
	return nil
}

func SendMail(msgs <-chan amqp.Delivery, c string) {
	for msg := range msgs {
		fmt.Println("收到消息", string(msg.Body))
		go Send(c, msg)
	}
}

func main() {
	var c *string
	c = flag.String("c", "", "消费者名称")
	flag.Parse()
	if *c == "" {
		log.Fatal("c参数一定要写")
	}
	mq := Lib.NewMQ()

	//处理消费者限流
	//连续发出两个消息，收到两个消息(ack)之后再接收后面的消息
	err := mq.Channel.Qos(2, 0, false) //最多能连续发两条消息，直到收到ack，才能继续发

	if err != nil {
		log.Fatal(err)
	}
	mq.Consume(Lib.QUEUE_NEWUSER, *c, SendMail)

	defer mq.Channel.Close()
}
