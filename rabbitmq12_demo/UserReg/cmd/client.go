package main

import (
	"Go_RabbitMQ/rabbitmq12_demo/AppInit"
	"Go_RabbitMQ/rabbitmq12_demo/Helper"
	"Go_RabbitMQ/rabbitmq12_demo/Lib"
	"flag"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

//假设这是真正的发送邮件的函数
func Send(c string, msg amqp.Delivery) {
	time.Sleep(time.Second * 1) //假设这是操作 邮件的调用 并假设操作很耗时
	userID := string(msg.Body)
	isfail := true                  //这个标志代表假定 失败了---一直是失败的
	delay := msg.Headers["x-delay"] //原来的延迟时间
	if isfail {
		r := Helper.SetNotify(userID, 5) //最大重试5次
		if r > 0 {
			newDelay := int(delay.(int32)) * 2 //每次收到消息 延迟时间*2 单位：毫秒
			err := myclient.SendDelayMessage(Lib.ROUTER_KEY_USERREG, Lib.EXCHANGE_USER_DELAY, userID, newDelay)
			if err != nil {
				log.Println(err)
			}
			log.Printf("%s向userID=%s的用户发送邮件:重试---延迟是:%d\n", c, string(msg.Body), newDelay)
		} else {
			log.Println("达到了最大次数，不再重发")
		}
		msg.Reject(false) //丢弃原消息
	} else { //假设邮件发送成功
		fmt.Printf("%s向userID=%s的用户发送邮件\n", c, string(msg.Body))
		msg.Ack(false)
	}
}

func SendMail(msgs <-chan amqp.Delivery, c string) {
	for msg := range msgs {
		fmt.Println("收到消息", string(msg.Body))
		go Send(c, msg)
	}

}

var myclient *Lib.MQ

func main() {
	var c *string
	c = flag.String("c", "", "消费者名称")
	flag.Parse()
	if *c == "" {
		log.Fatal("c参数一定要写")
	}
	dberr := AppInit.DBInit() //db初始化
	if dberr != nil {
		log.Fatal("DB error:", dberr)
	}
	mq := Lib.NewMQ()
	err := mq.Channel.Qos(2, 0, false) //最多能连续发两条消息，直到收到ack，才能继续发
	if err != nil {
		log.Fatal(err)
	}
	mq.Consume(Lib.QUEUE_NEWUSER, *c, SendMail)

	defer mq.Channel.Close()
}
