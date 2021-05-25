package main

import (
	"Go_RabbitMQ/rabbitmq07_demo/Lib"
	"Go_RabbitMQ/rabbitmq07_demo/UserReg/Models"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"

	"time"
)

func main() {
	router := gin.Default()
	router.Handle("POST", "/user", func(context *gin.Context) {
		userModel := Models.NewUserModel()
		err := context.BindJSON(&userModel)
		if err != nil {
			context.JSON(400, gin.H{
				"result": "param error",
			})
		} else {
			userModel.UserId = int(time.Now().Unix())
			if userModel.UserId > 0 {
				mq := Lib.NewMQ()
				err := mq.SendMessage(Lib.ROUTER_KEY_USERREG, Lib.EXCHANGE_USER, strconv.Itoa(userModel.UserId))
				defer mq.Channel.Close()
				if err != nil {
					log.Println(err)
				}
			}
			context.JSON(200, gin.H{
				"result": userModel,
			})
		}
	})
	c := make(chan error)
	go func() {
		err := router.Run(":8080")
		if err != nil {
			c <- err
		}
	}()

	go func() {
		err := Lib.UserInit() //初始化用户队列
		if err != nil {
			c <- err
		}
	}()
	err := <-c
	log.Fatal(err)
}
