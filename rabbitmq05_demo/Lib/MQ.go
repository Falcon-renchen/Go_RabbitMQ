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
	_, err := this.Channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}
	return this.Channel.Publish("", queueName, false, false,
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
