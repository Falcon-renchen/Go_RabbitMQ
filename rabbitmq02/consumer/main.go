package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"log"
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

	err = ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		fmt.Printf("ch.Qos() failed, err:%v\n", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true, //关闭消息自动确认
		false,
		false,
		false,
		nil,
	)

	//开启循环不断地消费消息，，，阻塞
	forever := make(chan bool)

	go func() {
		//for d := range msgs {
		//	log.Printf("Received a message: %s", d.Body)
		//	dot_count := bytes.Count(d.Body, []byte("."))  // 数一下有几个.  有几个.就有几个消息任务
		//	t := time.Duration(dot_count)
		//	time.Sleep(t * time.Second)  // 模拟耗时的任务
		//	log.Printf("Done")
		//}
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false) // 手动传递消息确认
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
