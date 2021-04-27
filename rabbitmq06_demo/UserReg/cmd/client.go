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
			//如果邮件没有被接收，则重启之后会自动接收。
			fmt.Printf("%s向userID=%s的用户发送邮件\n", c, string(msg.Body))
			//fmt.Println(msg.DeliveryTag, string(msg.Body))
			time.Sleep(time.Second * 1)
		}
		if c == "c1" { //模拟c1出了问题
			//如果没有这一行，c1会收到邮件，但是没有ack。如果把c1停掉，邮件会被c2收到
			msg.Reject(true) //重新入列，c1和c2都会收到消息	//  msg.Reject(false)   如果c1出现问题，丢弃消息。c1和c2都没有消息
			continue
		}

		//没有这段代码，说明c1，c2正常，先发给c1，然后发给c2，默认为轮询方式
		//if c == "c1" { //模拟c1出了问题
		//	msg.Reject(true) //重新入列		//c1出了问题，会重新入列
		//	continue
		//}
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
