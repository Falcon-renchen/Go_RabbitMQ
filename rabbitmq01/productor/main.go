package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func main() {
	fmt.Println("Go RabbitMQ")

	// 1. 尝试连接RabbitMQ，建立连接
	// 该连接抽象了套接字连接，并为我们处理协议版本协商和认证等。
	conn, err := amqp.Dial("amqp://guest:guest@172.16.17.152:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully Connected To our RabbitMQ Instance")


	// 2. 接下来，我们创建一个通道，大多数API都是用过该通道操作的。
	//处理大量的amqp，
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	// 3. 要发送，我们必须声明要发送到的队列。
	//在rabbitmq上声明一个队列，用来保存消息并且传递给消费者
	q, err := ch.QueueDeclare("TestQueue",false,false,false,false,nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//获取队列
	fmt.Println(q)

	// 4. 然后我们可以将消息发布到声明的队列
	//有效的发布来自客户端的消息
	err = ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			Headers:         nil,
			ContentType:     "text/plain",
			ContentEncoding: "",
			DeliveryMode:    0,
			Priority:        0,
			CorrelationId:   "",
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       "",
			Timestamp:       time.Time{},
			Type:            "",
			UserId:          "",
			AppId:           "",
			Body:            []byte("hello world"),
		},
		)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully Published Message to Queue")




}
