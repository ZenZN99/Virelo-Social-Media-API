package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func FollowRoutes(
	r *gin.Engine,
	followHandler *handlers.FollowHandler,
	authMiddleware gin.HandlerFunc,
) {
	api := r.Group("/api/follow")
	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.POST("/", followHandler.FollowUser)
	protected.DELETE("/:userId", followHandler.UnfollowUser)

	protected.GET("/is-following", followHandler.IsFollowing)
	protected.GET("/followers/:userId", followHandler.GetFollowers)
	protected.GET("/following/:userId", followHandler.GetFollowing)
	protected.GET("/counts/:userId", followHandler.GetCounts)
}
