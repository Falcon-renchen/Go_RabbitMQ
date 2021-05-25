package Lib

import (
	"Go_RabbitMQ/rabbitmq10_demo/AppInit"

	"github.com/streadway/amqp"
	"log"
	"strings"
	"time"
)

const (
	QUEUE_NEWUSER       = "newuser"           //用户注册 对应的队列名称
	QUEUE_NEWUSER_UNION = "newuser_union"     //合作单位用户注册 对应的队列名称
	EXCHANGE_USER       = "UserExchange"      //用户模块相关的交换机
	EXCHANGE_USER_DELAY = "UserExchangeDelay" //延迟队列
	ROUTER_KEY_USERREG  = "userreg"           //注册用户的路由key
)

type MQ struct {
	Channel       *amqp.Channel
	notifyConfirm chan amqp.Confirmation
	notifyReturn  chan amqp.Return
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

func (this *MQ) DecQueueAndBindWithArgs(name string, s string, name2 string, args map[string]interface{}) error {
	qList := strings.Split(name, ",")
	for _, queue := range qList {
		q, err := this.Channel.QueueDeclare(queue, false, false, false, false, args)
		if err != nil {
			return err
		}
		err = this.Channel.QueueBind(q.Name, s, name2, false, args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *MQ) NotifyReturn() {
	//如果消息没有正常进入队列，
	this.notifyReturn = this.Channel.NotifyReturn(make(chan amqp.Return))
	go this.listenReturn()
}

func (this *MQ) listenReturn() {
	ret := <-this.notifyReturn
	if string(ret.Body) != "" {
		log.Println("当前消息没有正确入列", string(ret.Body))
	}

}

//确认机制
func (this *MQ) SetConfirm() {
	err := this.Channel.Confirm(false)
	if err != nil {
		log.Fatal(err)
	}
	this.notifyConfirm = this.Channel.NotifyPublish(make(chan amqp.Confirmation))

}

func (this *MQ) ListenConfirm() {
	defer this.Channel.Close()
	ret := <-this.notifyConfirm
	if ret.Ack {
		log.Println("消息发送成功")
	} else {
		log.Println("消息发送失败")
	}
}

func (this *MQ) SendMessage(key string, exchange string, message string) error {
	err := this.Channel.Publish(exchange, key, true, false,
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
	return err
}

func (this *MQ) SendDelayMessage(key string, exchange string, message string, delay int) error {
	err := this.Channel.Publish(exchange, key, true, false,
		amqp.Publishing{
			Headers:         map[string]interface{}{"x-delay": delay},
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
	return err
}

func (this *MQ) Consume(queue string, key string, callback func(<-chan amqp.Delivery, string)) {
	msgs, err := this.Channel.Consume(queue, key, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	callback(msgs, key)

}
