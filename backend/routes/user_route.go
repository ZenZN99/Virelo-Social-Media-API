package routes

import (
	"backend/handlers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(
	r *gin.Engine,
	userHandler *handlers.UserHandler,
	authMiddleware gin.HandlerFunc,
	adminMiddleware gin.HandlerFunc,
) {

	api := r.Group("/api/auth")

	api.POST("/signup", userHandler.SignUp)
	api.POST("/login", userHandler.Login)

	protected := api.Group("/")
	protected.Use(authMiddleware)

	protected.GET("/me", userHandler.Me)
	protected.PUT("/update/profile", userHandler.UpdateProfile)
	protected.GET("/users", userHandler.GetAllUsers)
	protected.GET("/user/:userId", userHandler.GetUserById)

	protected.DELETE("/user/:userId",
		adminMiddleware,
		userHandler.DeleteUser,
	)
}
