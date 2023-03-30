package main

import (
	// "fmt"

	"fmt"
	"log"
	"net/http"
	"startup/auth"
	"startup/handler"
	"startup/helper"
	"startup/user"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/startup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepositoryUser(db)
	userService := user.NewService(userRepository)

	// Debug find email
	/*userByEmail, err  := userRepository.FindByEmail("addmin@gmail.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	if (userByEmail.ID == 0) {
		fmt.Println("user tidak ditentukan")
	} else {
		fmt.Println(userByEmail.Name)
	}
	fmt.Println(userByEmail)*/

	//testing password match compare
	/*	input := user.LoginInput{
			Email: "admin@gmail.com",
			Password: "admine",

		}
		user, err := userService.Login(input)
		if err != nil {
			fmt.Println("terjadi kesalahan")
			fmt.Println(err.Error())
		}
		fmt.Println(user.Email)
		fmt.Println(user.Name)*/

	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)
	// userService.SaveAvatar(1,"images/profile.png")testing langsung upload ke DB

	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2fQ.spLHRa9t4bBQ74eF8vYtWF-x8vc98NKXeWY-gwkdzJM")
	if err != nil {
		fmt.Println("ERROR")
		fmt.Println("ERROR")
		fmt.Println("ERROR")
	}
	if token.Valid {
		fmt.Println("VALID")
		fmt.Println("VALID")
	} else {
		fmt.Println("NOT VALID")
		fmt.Println("NOT VALID")
	}

	router := gin.Default()
	api := router.Group("api/v1")
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checker", userHandler.CheckEmailAvailability)
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	router.Run()

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		autHeader := c.GetHeader("Authorization")

		if !strings.Contains(autHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := ""
		// Bearer tokentoken
		arrayToken := strings.Split(autHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(claim["user_id"].(float64))
		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		
		c.Set("currentUser", user)
	}
}
