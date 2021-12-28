package main

import (
	"backend-example/auth"
	"backend-example/redis"
	"backend-example/user"
	"fmt"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	redis.SetupRedisClient()

	router = gin.Default()

	securedApiGroup := auth.SetupAuthentication(router)
	user.SetupRoutes(securedApiGroup)

	err := router.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
