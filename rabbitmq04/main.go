package main

import (
	"Go_RabbitMQ/rabbitmq04/AppInit"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn := AppInit.GetConn()
	defer conn.Close()

	//先创建一个Channel，分发
	c, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	//队列创建成功
	queue, err := c.QueueDeclare(
		"test",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("test003"),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("发送消息成功")
}
