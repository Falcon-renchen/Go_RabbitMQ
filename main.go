package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func main() {
	fmt.Println("Go RabbitMQ")

	conn, err := amqp.Dial("amqp://guest:guest@172.16.17.152:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully Connected To our RabbitMQ Instance")


	//处理大量的amqp，
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	//在rabbitmq上声明一个队列，用来保存消息并且传递给消费者
	q, err := ch.QueueDeclare("TestQueue",false,false,false,false,nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//获取队列
	fmt.Println(q)


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
