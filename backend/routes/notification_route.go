package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(r *gin.Engine, h *handlers.NotificationHandler, authMiddleware gin.HandlerFunc) {

	api := r.Group("/api/notification")
	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.GET("/", h.GetMyNotifications)
	protected.GET("/unread-count", h.GetUnreadCount)
	protected.PUT("/:id/read", h.MarkAsRead)
	protected.PATCH("/read-all", h.MarkAllAsRead)
	protected.DELETE("/:notificationId", h.DeleteNotification)
}
