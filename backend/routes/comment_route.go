package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.Engine, commentHandler *handlers.CommentHandler, authMiddleware gin.HandlerFunc) {
	api := r.Group("/api/comment")

	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.POST("/create", commentHandler.CreateComment)
	protected.PUT("/update/:commentId", commentHandler.UpdateComment)
	protected.DELETE("/delete/:commentId", commentHandler.DeleteComment)
	protected.GET("/comments", commentHandler.GetComments)
	protected.GET("/user/:userId", commentHandler.GetCommentsByUser)
	protected.GET("/replies/:commentId", commentHandler.GetReplies)
}
