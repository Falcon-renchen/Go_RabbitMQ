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
	//qs := fmt.Sprintf("%s", QUEUE_NEWUSER)
	err = mq.DecQueueAndBind(qs, ROUTER_KEY_USERREG, EXCHANGE_USER)
	if err != nil {
		return fmt.Errorf("Queue Bind error", err)
	}
	return nil
}

//延迟队列
func UserDelayInit() error {
	mq := NewMQ()
	if mq == nil {
		return fmt.Errorf("mq init error")
	}
	defer mq.Channel.Close()

	//申明交换机
	err := mq.Channel.ExchangeDeclare(EXCHANGE_USER_DELAY, "direct",
		false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("DelayExchange error", err)
	}
	//qs := fmt.Sprintf("%s", QUEUE_NEWUSER)
	qs := fmt.Sprintf("%s", QUEUE_NEWUSER)
	args := map[string]interface{}{}
	err = mq.DecQueueAndBindWithArgs(qs, ROUTER_KEY_USERREG, EXCHANGE_USER_DELAY, args)
	if err != nil {
		return fmt.Errorf("Delay Queue Bind error", err)
	}
	return nil
}
