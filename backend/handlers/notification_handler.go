package handlers

import (
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// GET /api/notification
func (h *NotificationHandler) GetMyNotifications(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	notifications, err := h.notificationService.GetUserNotifications(
		c.Request.Context(),
		payload.UserID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GET /api/notification/unread-count
func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	count, err := h.notificationService.GetUnreadCount(
		c.Request.Context(),
		payload.UserID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"unreadCount": count,
	})
}

// PUT /api/notification/:id/read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {

	id := c.Param("id")

	notification, err := h.notificationService.MarkAsRead(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notification)
}

// PATCH /api/notification/read-all
func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {

	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	payload := userData.(*utils.TokenPayload)

	res, err := h.notificationService.MarkAllAsRead(
		c.Request.Context(),
		payload.UserID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// DELETE /api/notification/:notificationId
func (h *NotificationHandler) DeleteNotification(c *gin.Context) {

	id := c.Param("notificationId")

	res, err := h.notificationService.DeleteNotification(
		c.Request.Context(),
		id,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
