package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Consumer Application")
	conn, err := amqp.Dial("amqp://wyp:123@172.16.17.154:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"TestQueue",
		"c1",
		true,
		false,
		false,
		false,
		nil,
	)

	//处理消息
	//阻塞，直到收到一个消息，不然程序会直接结束
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Printf("Receving Message: %s\n", d.Body)
		}
	}()
	fmt.Println("Successfully conntected to our rabbitmq instance")
	fmt.Println("[*] -- waitting for messages")

	<-forever //在终端中退出

}
