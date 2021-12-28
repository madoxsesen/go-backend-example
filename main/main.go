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

	auth.SetupRoutes(router)
	user.SetupRoutes(router)

	auth.SetupMiddleware(router)

	err := router.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
