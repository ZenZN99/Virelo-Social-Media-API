package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) SignUp(c *gin.Context) {

	var data models.UserModel

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.userService.SignUp(c.Request.Context(), data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) Login(c *gin.Context) {

	var data models.UserModel

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	res, err := h.userService.Login(c.Request.Context(), data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Me(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	user, err := h.userService.Me(c.Request.Context(), payload.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	bio := c.PostForm("bio")

	file, err := c.FormFile("avatar")
	var filePath string

	if err == nil {
		filePath = "./uploads/" + file.Filename

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to upload file",
			})
			return
		}
	}

	res, err := h.userService.UpdateProfile(
		c.Request.Context(),
		payload.UserID,
		filePath,
		bio,
	)

	if err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	users, err := h.userService.GetAllUsers(
		c.Request.Context(),
		payload.UserID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func (h *UserHandler) GetUserById(c *gin.Context) {

	userId := c.Param("userId")

	user, err := h.userService.GetUserById(
		c.Request.Context(),
		userId,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {

	userId := c.Param("userId")

	_, err := h.userService.DeleteUserById(
		c.Request.Context(),
		userId,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
