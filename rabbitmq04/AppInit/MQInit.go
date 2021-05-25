package AppInit

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

var MQConn *amqp.Connection

func init() {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d", "wyp", "123", "172.16.17.155", 5672)
	conn, err := amqp.Dial(dsn)
	if err != nil {
		log.Fatal(err)
	}
	MQConn = conn
	log.Println(MQConn.Major)

}

func GetConn() *amqp.Connection {
	return MQConn
}
