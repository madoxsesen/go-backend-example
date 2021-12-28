package user

import (
	"backend-example/auth"
	"backend-example/redis"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

const routePrefix = "/users"

type Data struct {
	UserInfo int `json:"userInfo"`
}

func (u *Data) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func SetupRoutes(router *gin.RouterGroup) {
	router.GET(routePrefix, getUserInformation)
	router.PUT(routePrefix, editUserInformation)
}

func getUserInformation(c *gin.Context) {
	res, _ := c.Get("email")
	token := res.(*auth.JwtToken)

	fmt.Println(token)

	foo, _ := redis.Client.Get(redis.Client.Context(), token.Email+"-data").Result()

	var userData Data

	err := json.Unmarshal([]byte(foo), &userData)
	if err != nil {
		return
	}

	fmt.Println(userData.UserInfo)
	c.JSON(200, &userData)
}

func editUserInformation(c *gin.Context) {
	res, _ := c.Get("email")
	token := res.(*auth.JwtToken)

	result, err := redis.Client.Set(redis.Client.Context(), token.Email+"-data", &Data{UserInfo: 1234}, 0).Result()
	if err != nil {
		return
	}

	fmt.Println(result)
	c.Status(200)
}
