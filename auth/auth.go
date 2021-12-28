package auth

import (
	"backend-example/redis"
	"encoding/json"
	"fmt"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type AuthenticationDTO struct {
	Email    string
	Password string
}

type UserInformation struct {
	Email        string
	PasswordHash string
}

func (u *UserInformation) String() string {
	return u.Email + " " + u.PasswordHash
}

func (u *UserInformation) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
}

func register(c *gin.Context) {
	var registrationDTO AuthenticationDTO

	err := c.BindJSON(&registrationDTO)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Printf("JwtToken: %s, Pass: %s\n", registrationDTO.Email, registrationDTO.Password)

	hashedPassword, err := hashPassword(registrationDTO.Password)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	var userInformation = UserInformation{Email: registrationDTO.Email, PasswordHash: string(hashedPassword)}

	result, err := redis.Client.Set(redis.Client.Context(), registrationDTO.Email, &userInformation, 0).Result()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Printf("Result: %s\n", result)
}

func hashPassword(pw string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pw), 0)
}

func login(c *gin.Context) (interface{}, error) {
	var loginDTO AuthenticationDTO

	err := c.BindJSON(&loginDTO)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return nil, err
	}

	fmt.Printf("JwtToken: %s, Pass: %s\n", loginDTO.Email, loginDTO.Password)

	jsonMarshalled, err := redis.Client.Get(redis.Client.Context(), loginDTO.Email).Result()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return nil, err
	}

	fmt.Printf("Result: %s\n", jsonMarshalled)

	var userInformation UserInformation

	err = json.Unmarshal([]byte(jsonMarshalled), &userInformation)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return nil, err
	}

	fmt.Printf("Parsed: %s\n", userInformation)

	err = bcrypt.CompareHashAndPassword([]byte(userInformation.PasswordHash), []byte(loginDTO.Password))
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return nil, err
	}

	fmt.Println("Login succeeded!")

	return &JwtToken{Email: userInformation.Email}, nil
}

var identityKey = "email"

type JwtToken struct {
	Email string
}

func SetupAuthentication(c *gin.Engine) *gin.RouterGroup {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "backendExample",
		Key:         []byte("highlysupersecretkey"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*JwtToken); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &JwtToken{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: login,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(*JwtToken); ok {
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

	return setupRoutes(authMiddleware, c)
}

func setupRoutes(authMiddleware *jwt.GinJWTMiddleware, c *gin.Engine) *gin.RouterGroup {
	errInit := authMiddleware.MiddlewareInit()

	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// Register outside group, so it does not use middleware
	c.POST("/api/auth/login", authMiddleware.LoginHandler)
	c.POST("/api/auth/register", register)

	auth := c.Group("/api")
	auth.GET("/auth/refresh", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())

	return auth
}
