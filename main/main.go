package main

import (
	"backend-example/auth"
	"backend-example/user"
	"fmt"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	router = gin.Default()

	auth.SetupRoutes(router)
	user.SetupRoutes(router)

	err := router.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
}
