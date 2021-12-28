package auth

import "github.com/gin-gonic/gin"

const routePrefix = "/api/auth/"

func SetupRoutes(router *gin.Engine) {
	router.POST(routePrefix+"register", register)
	router.POST(routePrefix+"login", login)
}

func register(c *gin.Context) {

}

func login(c *gin.Context) {

}
