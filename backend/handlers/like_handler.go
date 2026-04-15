package handlers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeService *services.LikeService
}

func NewLikeHandler(likeService *services.LikeService) *LikeHandler {
	return &LikeHandler{
		likeService: likeService,
	}
}

func (h *LikeHandler) ToggleLike(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	var body struct {
		TargetID   string          `json:"targetId"`
		TargetType models.ContentType `json:"targetType"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	res, err := h.likeService.ToggleLike(
		c.Request.Context(),
		payload.UserID,
		body.TargetID,
		body.TargetType,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *LikeHandler) CountLikes(c *gin.Context) {

	targetID := c.Query("targetId")
	targetType := models.ContentType(c.Query("targetType"))

	count, err := h.likeService.CountLikes(
		c.Request.Context(),
		targetID,
		targetType,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

func (h *LikeHandler) IsLiked(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	targetID := c.Query("targetId")
	targetType := models.ContentType(c.Query("targetType"))

	liked, err := h.likeService.IsLiked(
		c.Request.Context(),
		payload.UserID,
		targetID,
		targetType,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"liked": liked,
	})
}

func (h *LikeHandler) GetLikes(c *gin.Context) {

	targetID := c.Query("targetId")
	targetType := models.ContentType(c.Query("targetType"))

	likes, err := h.likeService.GetLikes(
		c.Request.Context(),
		targetID,
		targetType,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"likes": likes,
	})
}
