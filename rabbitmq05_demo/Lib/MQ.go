package Lib

import (
	"Go_RabbitMQ/rabbitmq05_demo/AppInit"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const (
	QUEUE_NEWUSER = "newuser" //用户注册 对应的队列名称
)

type MQ struct {
	Channel *amqp.Channel
}

func NewMQ() *MQ {
	c, err := AppInit.GetConn().Channel()
	if err != nil {
		log.Println(err)
		return nil
	}
	return &MQ{Channel: c}
}

func (this *MQ) SendMessage(queueName string, message string) error {
	q1, err := this.Channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	//假设是合作网站用的队列
	q2, err := this.Channel.QueueDeclare(queueName+"union", false, false, false, false, nil)
	if err != nil {
		return err
	}
	//声明交换机
	err = this.Channel.ExchangeDeclare("UserExchange", "direct", false, false, false, false, nil)
	if err != nil {
		return err
	}
	//将q1.Name 绑定到userreg路由键，将所有关联到UserExchange里面
	err = this.Channel.QueueBind(q1.Name, "userreg", "UserExchange", false, nil)
	if err != nil {
		return err
	}
	err = this.Channel.QueueBind(q2.Name, "userreg", "UserExchange", false, nil)
	if err != nil {
		return err
	}

	return this.Channel.Publish("UserExchange", "userreg", false, false,
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
			Body:            []byte(message),
		})
}
