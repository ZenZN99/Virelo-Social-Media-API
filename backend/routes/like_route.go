package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func LikeRoutes(r *gin.Engine, likeHandler *handlers.LikeHandler, authMiddleware gin.HandlerFunc) {

	api := r.Group("/api/likes")
	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.POST("/toggle", likeHandler.ToggleLike)
	protected.GET("/count", likeHandler.CountLikes)
	protected.GET("/is-liked", likeHandler.IsLiked)
	protected.GET("/list", likeHandler.GetLikes)

}
