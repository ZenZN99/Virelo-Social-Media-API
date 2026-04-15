package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func PostRoutes(
	r *gin.Engine,
	postHandler *handlers.PostHandler,
	authMiddleware gin.HandlerFunc,
) {
	api := r.Group("/api/post")

	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.POST("/create", postHandler.CreatePost)
	protected.PUT("/update/:postId", postHandler.UpdatePost)
	protected.DELETE("/delete/:postId", postHandler.DeletePost)
	protected.GET("/posts", postHandler.GetAllPosts)
	protected.GET("/post/:postId", postHandler.GetPostByID)
	protected.GET("/user", postHandler.GetPostsByUser)
}
