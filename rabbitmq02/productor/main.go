package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"strings"
	"time"
)

func main() {
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
	fmt.Println(q)


	body := bodyForm(os.Args)
	err = ch.Publish("","TestQueue",false,false,amqp.Publishing{
		Headers:         nil,
		ContentType:     "text/plain",
		ContentEncoding: "",
		DeliveryMode:    amqp.Persistent,  // 持久（交付模式：瞬态/持久）
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       "",
		Timestamp:       time.Time{},
		Type:            "",
		UserId:          "",
		AppId:           "",
		Body:            []byte(body),
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully Published Message to Queue")
}



func bodyForm(args []string) string {
	var s string

	if len(args)<2 || os.Args[1]=="" {
		s = "hello"
	} else {
		s = strings.Join(args[1:],"")
	}
	return s
}