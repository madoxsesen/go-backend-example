package auth

import (
	"backend-example/redis"
	"encoding/json"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var identityKey = "email"

type User struct {
	Email string
}

func SetupMiddleware(c *gin.Engine) *gin.RouterGroup {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "backendExample",
		Key:         []byte("highlysupersecretkey"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginDTO AuthenticationDTO

			err := c.BindJSON(&loginDTO)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				return nil, nil
			}

			fmt.Printf("User: %s, Pass: %s\n", loginDTO.Email, loginDTO.Password)

			jsonMarshalled, err := redis.Client.Get(redis.Client.Context(), loginDTO.Email).Result()
			if err != nil {
				fmt.Println("Error: " + err.Error())
				return nil, nil
			}

			fmt.Printf("Result: %s\n", jsonMarshalled)

			var userInformation UserInformation

			err = json.Unmarshal([]byte(jsonMarshalled), &userInformation)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				return nil, nil
			}

			fmt.Printf("Parsed: %s\n", userInformation)

			err = bcrypt.CompareHashAndPassword([]byte(userInformation.PasswordHash), []byte(loginDTO.Password))
			if err != nil {
				fmt.Println("Error: " + err.Error())
				return nil, nil
			}

			fmt.Println("Login succeeded!")
			return &User{Email: loginDTO.Email}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*User); ok {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	c.POST("/jwtlogin", authMiddleware.LoginHandler)
	c.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := c.Group("/api")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())

	return auth
}
