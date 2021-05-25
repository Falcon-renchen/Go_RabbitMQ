package Lib

import (
	"Go_RabbitMQ/rabbitmq08_demo/AppInit"

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

//监听消息入列回值
func (this *MQ) NotifyReturn() {
	//如果消息没有正常进入队列，会给notifychannel插入一个值
	this.notifyReturn = this.Channel.NotifyReturn(make(chan amqp.Return))
	go this.listenReturn()
}

func (this *MQ) listenReturn() {
	ret := <-this.notifyReturn
	//会返回消息体， "消息没有正常入列 1234556" "confirm消息发送成功"
	//exchange删掉会出现 消息发送失败
	//点击会正常，解绑或者queue删掉，会出现消息没有正常入列
	if string(ret.Body) != "" {
		log.Println("当前消息没有正确入列", string(ret.Body))
	}

}

//确认机制，MQ会出现错误，出现回值，做一个日志来记录，然后重发或者手动发
func (this *MQ) SetConfirm() {
	//设置confirm模式
	err := this.Channel.Confirm(false)
	if err != nil {
		log.Fatal(err)
	}
	this.notifyConfirm = this.Channel.NotifyPublish(make(chan amqp.Confirmation))
}

//监听确认机制返回的值。
func (this *MQ) ListenConfirm() {
	defer this.Channel.Close()
	ret := <-this.notifyConfirm
	if ret.Ack {
		log.Println("confirm：消息发送成功")
	} else {
		log.Println("confirm：消息发送失败")
	}
}

func (this *MQ) SendMessage(key string, exchange string, message string) error {
	//mandatory 为 true，在exchange正常且可达的情况，
	//如果exchange+routekey无法传递给queue，mq会将消息返还给消费者
	//如果为false，则直接丢弃
	return this.Channel.Publish(exchange, key, true, false,
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

func (this *MQ) Consume(queue string, key string, callback func(<-chan amqp.Delivery, string)) {
	msgs, err := this.Channel.Consume(queue, key, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	callback(msgs, key)

}
