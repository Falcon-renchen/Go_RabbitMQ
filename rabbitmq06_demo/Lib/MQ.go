package Lib

import (
	"Go_RabbitMQ/rabbitmq06_demo/AppInit"

	"github.com/streadway/amqp"
	"log"
	"strings"
	"time"
)

const (
	QUEUE_NEWUSER       = "newuser"       //用户注册 对应的队列名称
	QUEUE_NEWUSER_UNION = "newuser_union" //合作单位用户注册 对应的队列名称
	EXCHANGE_USER       = "UserExchange"  //用户模块相关的交换机
	ROUTER_KEY_USERREG  = "userreg"       //注册用户的路由key
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

//申明队列以及绑定路由key
//多个队列 可以用逗号分割
func (this *MQ) DecQueueAndBind(queues string, key string, exchange string) error {
	qList := strings.Split(queues, ",")
	for _, queue := range qList {
		q, err := this.Channel.QueueDeclare(queue, false, false, false, false, nil)
		if err != nil {
			return err
		}
		err = this.Channel.QueueBind(q.Name, key, exchange, false, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MQ) SendMessage(key string, exchange string, message string) error {
	return this.Channel.Publish(exchange, key, false, false,
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
		},
	)
}

//key 是消费者名称
func (this *MQ) Consume(queue string, key string, callback func(<-chan amqp.Delivery, string)) {
	msgs, err := this.Channel.Consume(queue, key, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	callback(msgs, key)

}
