package Lib

import "fmt"

//初始化用户相关的队列
func UserInit() error {
	mq := NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	//申明交换机
	err := mq.Channel.ExchangeDeclare(EXCHANGE_USER, "direct", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("Exchange error")
	}
	qs := fmt.Sprintf("%s,%s", QUEUE_NEWUSER, QUEUE_NEWUSER_UNION)
	err = mq.DecQueueAndBind(qs, ROUTER_KEY_USERREG, EXCHANGE_USER)
	if err != nil {
		return fmt.Errorf("Queue Bind error", err)
	}
	return nil
}

func UserDelayInit() error {
	mq := NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	//申明交换机
	err := mq.Channel.ExchangeDeclare(EXCHANGE_USER_DELAY, "x-delayed-message",
		false, false, false, false,
		map[string]interface{}{"x-delayed-type": "direct"})
	if err != nil {
		return fmt.Errorf("Exchange error")
	}
	qs := fmt.Sprintf("%s,%s", QUEUE_NEWUSER, QUEUE_NEWUSER_UNION)
	err = mq.DecQueueAndBind(qs, ROUTER_KEY_USERREG, EXCHANGE_USER)
	if err != nil {
		return fmt.Errorf("Delay Queue Bind error", err)
	}
	return nil
}
