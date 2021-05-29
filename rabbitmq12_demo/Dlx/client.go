package main

import (
	"Go_RabbitMQ/rabbitmq12_demo/Lib"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func TestDlx(msgs <-chan amqp.Delivery, c string) {
	for msg := range msgs {
		fmt.Println("收到消息", string(msg.Body))
		msg.Ack(false)
	}
}
func main() {
	client := Lib.NewMQ()
	err := Lib.DlXTestInit()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Channel.Qos(2, 0, false)
	if err != nil {
		log.Fatal(err)
	}
	client.Consume("test_dlx", "c", TestDlx)

	defer client.Channel.Close()
}
