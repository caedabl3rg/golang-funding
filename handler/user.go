package handler

import (
	"fmt"
	"net/http"
	"startup/auth"
	"startup/helper"
	"startup/user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
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
	token, err := h.authService.GenerateToken(newUser.ID)
	formatter := user.FormatUser(newUser, token)
	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) Login(c *gin.Context) {
	// USER MEMASUKAN INPUT EMAIL DAN PASSWORD
	// INPUT DITANGKAP HANDLER
	// MAPPING DARI INPUT USER KE INPUT STRUCT
	// 2. INPUT STRUCT PASSING SERVICE
	// 1. DI SERVICE MENCARI DG BANTUAN REPOSITORY USER DENGAN EMAIL X
	// MENCOCOKAN PASSWORD

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormattError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	loggedUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login Failed", http.StatusUnprocessableEntity, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	token, err := h.authService.GenerateToken(loggedUser.ID)
	formatter := user.FormatUser(loggedUser, token)
	response := helper.APIResponse("Login Success", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	// ADA INPUT EMAIL DARI USER
	// INPUT EMAIL DI MAPPING KE STRUCT INPUT
	// STRUCT INPUT DI PASSING KE SERVICE
	// SERVICE AKAN MANGGIL REPOSITORY - EMAIL SUDAH ADA ATAU BELUM
	// REPOSITORY AKAN MELAKUKAN QUERY KE DB

	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormattError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Checking email failed", http.StatusUnprocessableEntity, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Errors"}
		response := helper.APIResponse("Checking email failed", http.StatusUnprocessableEntity, "errors", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	data := gin.H{
		"is_available": isEmailAvailable,
	}
	var metaMessage string
	if isEmailAvailable {
		metaMessage = "Email is available"
	} else {
		metaMessage = "Email has been registered"
	}
	response := helper.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// TANGKAP INPUT DARI USER FORM-BODY
	// SIMPAN GAMBAR DI FOLDER "images"
	// DISERVICE KITA PANGGIL REPO (UNTUK MENENTUKAN SIAPA USER YANG AKSES)
	// JWT
	// REPO AMBIL DATA USER YG ID = X
	// REPO UPDATE DATA USER SIMPAN LOKASI FILE

	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "errors", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// HARUS AMBIL DARI JWT
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	// path := "images/" + file.Filename
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "errors", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload avatar image", http.StatusBadRequest, "errors", data)
		c.JSON(http.StatusBadRequest, response)
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar success", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}
