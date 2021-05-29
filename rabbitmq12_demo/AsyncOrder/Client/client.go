package main

import (
	"Go_RabbitMQ/rabbitmq12_demo/AppInit"
	"Go_RabbitMQ/rabbitmq12_demo/AsyncOrder/models"
	"Go_RabbitMQ/rabbitmq12_demo/Lib"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func TestOrder(msgs <-chan amqp.Delivery, c string) {
	for msg := range msgs {
		//fmt.Println("收到消息",string(msg.Body))
		req := &models.RequestModel{}
		err := json.Unmarshal(msg.Body, req)
		if err != nil {
			log.Println(err)

		} else {
			fmt.Println(req)
			time.Sleep(time.Second * 2)
			_, err = AppInit.GetDB().NamedExec("insert into orders(order_no,order_time,order_user) "+
				"values(:orderno,:ordertime,:userid)", req)
			if err != nil {
				log.Println(err)
			}
		}

		msg.Ack(false)
	}
}
func main() {
	client := Lib.NewMQ()
	err := AppInit.DBInit()
	if err != nil {
		log.Fatal(err)
	}
	err = client.Channel.Qos(2, 0, false)
	if err != nil {
		log.Fatal(err)
	}
	client.Consume("Queue_Order", "c", TestOrder)

	defer client.Channel.Close()
}
