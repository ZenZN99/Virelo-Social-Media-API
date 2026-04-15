package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func MessageRoutes(
	r *gin.Engine,
	messageHandler *handlers.MessageHandler,
	authMiddleware gin.HandlerFunc,
) {

	api := r.Group("/api/message")

	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.POST("/send", messageHandler.SendMessage)

	protected.GET("/:receiverId", messageHandler.GetChatMessages)

	protected.DELETE("/:messageId", messageHandler.DeleteMessage)

	protected.PUT("/read/:senderId", messageHandler.MarkMessageAsRead)
}
