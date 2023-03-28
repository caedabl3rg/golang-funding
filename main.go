package main

import (
	// "fmt"
	"log"
	"startup/handler"
	"startup/user"

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
		
		
		
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("api/v1")

	// TESTING LOGIN
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checker", userHandler.CheckEmailAvailability)
	api.POST("/users", userHandler.RegisterUser)
	router.Run()

	// userInput := user.RegisterUserInput{}
	// userInput.Name = "Test simpan dari service"
	// userInput.Email = "contoh@gmail.com"
	// userInput.Occupation = "programmer"
	// userInput.Password = "password"
	// userService.RegisterUser(userInput)

}
