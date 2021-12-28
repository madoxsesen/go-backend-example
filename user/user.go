package user

import "github.com/gin-gonic/gin"

const routePrefix = "/api/users/"

func SetupRoutes(router *gin.Engine) {
	router.GET(routePrefix+":userId", getUserInformation)
	router.PUT(routePrefix+":userId", editUserInformation)
}

func getUserInformation(c *gin.Context) {

}

func editUserInformation(c *gin.Context) {

}
