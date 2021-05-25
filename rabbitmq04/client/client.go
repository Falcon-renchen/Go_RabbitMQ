package main

import (
	"Go_RabbitMQ/rabbitmq04/AppInit"
	"fmt"
	"log"
)

func main() {
	conn := AppInit.GetConn()
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	//go channel类型
	msgs, err := c.Consume(
		"test",
		"c1",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	//没有确认机制，，停掉之后会从unacked跳到ready，给其他消费者读取
	for msg := range msgs {
		fmt.Println(msg.DeliveryTag, string(msg.Body))
	}
}
