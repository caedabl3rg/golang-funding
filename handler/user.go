package handler

import (
	"net/http"
	"startup/helper"
	"startup/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func (h userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai parameter

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormattError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	//token, err := h.jwtService.Generate
	formatter := user.FormatUser(newUser, "tokentokentoken")
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService: userService}
}
