package handlers

import (
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService *services.FollowService
}

func NewFollowHandler(followService *services.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

func (h *FollowHandler) FollowUser(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	var body struct {
		FollowingID string `json:"followingId"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	follow, err := h.followService.FollowUser(
		c.Request.Context(),
		payload.UserID,
		body.FollowingID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "followed successfully",
		"data":    follow,
	})
}

func (h *FollowHandler) UnfollowUser(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	followingID := c.Param("userId")

	err := h.followService.UnfollowUser(
		c.Request.Context(),
		payload.UserID,
		followingID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "unfollowed successfully",
	})
}

func (h *FollowHandler) IsFollowing(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	targetID := c.Query("userId")

	isFollowing, err := h.followService.IsFollowing(
		c.Request.Context(),
		payload.UserID,
		targetID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"isFollowing": isFollowing,
	})
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {

	userID := c.Param("userId")

	followers, err := h.followService.GetFollowers(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
	})
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {

	userID := c.Param("userId")

	following, err := h.followService.GetFollowing(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"following": following,
	})
}

func (h *FollowHandler) GetCounts(c *gin.Context) {

	userID := c.Param("userId")

	followers, err := h.followService.GetFollowersCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	following, err := h.followService.GetFollowingCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"following": following,
	})
}
