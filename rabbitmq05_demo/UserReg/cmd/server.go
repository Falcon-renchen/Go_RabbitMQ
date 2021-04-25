package main

import (
	"Go_RabbitMQ/rabbitmq05_demo/Lib"
	"Go_RabbitMQ/rabbitmq05_demo/UserReg/Models"
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
			userModel.UserId = int(time.Now().Unix()) //假设是入库过程
			if userModel.UserId > 0 {                 //假设入库成功
				mq := Lib.NewMQ()
				err := mq.SendMessage(Lib.QUEUE_NEWUSER, strconv.Itoa(userModel.UserId))
				if err != nil {
					log.Println(err)
				}
			}
			context.JSON(200, gin.H{
				"result": userModel,
			})
		}
	})
	router.Run(":8080")
}
