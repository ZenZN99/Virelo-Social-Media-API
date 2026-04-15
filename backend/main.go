package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/app"
	"backend/routes"
	"backend/utils"
	"backend/websocket"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	container, err := app.NewContainer()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	utils.InitRedis()

	hub := websocket.NewHub()

	r.GET("/ws/chat", websocket.NewChatSocket(hub))

	routes.UserRoutes(
		r,
		container.Handlers.UserHandler,
		container.Middlewares.AuthMiddleware(),
		container.Middlewares.AdminMiddleware(),
	)

	routes.PostRoutes(
		r,
		container.Handlers.PostHandler,
		container.Middlewares.AuthMiddleware(),
	)

	routes.ReelRoutes(
		r,
		container.Handlers.ReelHandler,
		container.Middlewares.AuthMiddleware(),
	)

	routes.CommentRoutes(
		r,
		container.Handlers.CommentHandler,
		container.Middlewares.AuthMiddleware(),
	)

	routes.LikeRoutes(
		r,
		container.Handlers.LikeHandler,
		container.Middlewares.AuthMiddleware(),
	)

	routes.FollowRoutes(
		r,
		container.Handlers.FollowHandler,
		container.Middlewares.AuthMiddleware(),
	)

	routes.MessageRoutes(
		r,
		container.Handlers.MessageHandler,
		container.Middlewares.AuthMiddleware(),
	)

	routes.NotificationRoutes(
		r,
		container.Handlers.NotificationHandler,
		container.Middlewares.AuthMiddleware(),
	)

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Hello Golang"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on http://localhost:" + port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
