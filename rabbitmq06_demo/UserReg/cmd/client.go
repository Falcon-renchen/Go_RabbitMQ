package main

import (
	"Go_RabbitMQ/rabbitmq06_demo/Lib"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func SendMail(msgs <-chan amqp.Delivery, c string) {
	for msg := range msgs {
		{
			//模拟发送邮件
			fmt.Printf("%s向userID=%s的用户发送邮件\n", c, string(msg.Body))
			//fmt.Println(msg.DeliveryTag, string(msg.Body))
			time.Sleep(time.Second * 1)
		}
		if c == "c1" { //模拟c1出了问题
			msg.Reject(true) //重新入列
			continue
		}
		msg.Ack(false)
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
	mq.Consume(Lib.QUEUE_NEWUSER, *c, SendMail)

	defer mq.Channel.Close()
}
