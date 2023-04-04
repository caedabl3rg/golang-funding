package main

import (
	"log"
	"net/http"
	"startup/auth"
	"startup/campaign"
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
	// REPOSITORY
	userRepository := user.NewRepositoryUser(db)
	campaignRepository := campaign.NewRepositoryCampaign(db)

	// SERVICE
	userService := user.NewService(userRepository)
	campaignService := campaign.NewServiceCampaign(campaignRepository)

	// HANDLER
	authService := auth.NewService()
	userHandler := handler.NewUserHandler(userService, authService)
	campignHandler := handler.NewCampignHandler(campaignService)

	router := gin.Default()
	router.Static("/images","./images")
	api := router.Group("api/v1")
	api.POST("/sessions", userHandler.Login)
	api.POST("/email-checker", userHandler.CheckEmailAvailability)
	api.POST("/users", userHandler.RegisterUser)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)
	
	api.GET("/campaigns", campignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campignHandler.GetCampaign)



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
