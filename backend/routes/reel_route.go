package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func ReelRoutes(
	r *gin.Engine,
	reelHandler *handlers.ReelHandler,
	authMiddleware gin.HandlerFunc,
) {
	api := r.Group("/api/reel")

	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.POST("/create", reelHandler.CreateReel)
	protected.PUT("/update/:reelId", reelHandler.UpdateReel)
	protected.DELETE("/delete/:reelId", reelHandler.DeleteReel)
	protected.GET("/reels", reelHandler.GetAllReels)
	protected.GET("/reel/:reelId", reelHandler.GetReelById)
	protected.GET("/user", reelHandler.GetReelsByUser)
}
