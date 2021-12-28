package auth

import (
	"backend-example/redis"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const routePrefix = "/api/auth/"

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

func SetupRoutes(router *gin.Engine) {
	router.POST(routePrefix+"register", register)
	router.POST(routePrefix+"login", login)
}

func register(c *gin.Context) {
	var registrationDTO AuthenticationDTO

	err := c.BindJSON(&registrationDTO)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Printf("User: %s, Pass: %s\n", registrationDTO.Email, registrationDTO.Password)

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

func login(c *gin.Context) {
	var loginDTO AuthenticationDTO

	err := c.BindJSON(&loginDTO)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Printf("User: %s, Pass: %s\n", loginDTO.Email, loginDTO.Password)

	jsonMarshalled, err := redis.Client.Get(redis.Client.Context(), loginDTO.Email).Result()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Printf("Result: %s\n", jsonMarshalled)

	var userInformation UserInformation

	err = json.Unmarshal([]byte(jsonMarshalled), &userInformation)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Printf("Parsed: %s\n", userInformation)

	err = bcrypt.CompareHashAndPassword([]byte(userInformation.PasswordHash), []byte(loginDTO.Password))
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	fmt.Println("Login succeeded!")
}
