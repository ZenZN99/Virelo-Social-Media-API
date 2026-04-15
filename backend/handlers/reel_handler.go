package handlers

import (
	"backend/services"
	"backend/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReelHandler struct {
	reelService *services.ReelService
}

func NewReelHandler(reelService *services.ReelService) *ReelHandler {
	return &ReelHandler{
		reelService: reelService,
	}
}

func (h *ReelHandler) CreateReel(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	content := c.PostForm("content")

	file, err := c.FormFile("reel")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "reel file is required"})
		return
	}

	path := "/tmp/" + file.Filename

	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
		return
	}

	reel, err := h.reelService.CreateReel(
		c.Request.Context(),
		payload.UserID,
		content,
		path,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reel)
}

func (h *ReelHandler) UpdateReel(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	reelID := c.Param("reelId")

	var body struct {
		Content string `json:"content"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	reel, err := h.reelService.UpdateReel(
		c.Request.Context(),
		reelID,
		payload,
		body.Content,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reel)
}

func (h *ReelHandler) DeleteReel(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	reelID := c.Param("reelId")

	res, err := h.reelService.DeleteReel(
		c.Request.Context(),
		reelID,
		payload,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *ReelHandler) GetAllReels(c *gin.Context) {

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	reels, err := h.reelService.GetAllReels(
		c.Request.Context(),
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  reels,
	})
}

func (h *ReelHandler) GetReelById(c *gin.Context) {

	reelID := c.Param("reelId")

	reel, err := h.reelService.GetReelById(
		c.Request.Context(),
		reelID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reel)
}

func (h *ReelHandler) GetReelsByUser(c *gin.Context) {

	userID := c.Param("userId")

	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	reels, err := h.reelService.GetReelsByUser(
		c.Request.Context(),
		userID,
		page,
		limit,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"data":  reels,
	})
}
