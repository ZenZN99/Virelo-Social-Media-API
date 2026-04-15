package app

import (
	"backend/middlewares"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	AuthMiddleware  func() gin.HandlerFunc
	AdminMiddleware func() gin.HandlerFunc
}

func InitMiddlewares(tokenService *utils.TokenService) *Middlewares {
	authMiddleware := middlewares.AuthMiddleware(tokenService)
	adminMiddleware := middlewares.AdminMiddleware()

	return &Middlewares{
		AuthMiddleware: func() gin.HandlerFunc {
			return authMiddleware
		},
		AdminMiddleware: func() gin.HandlerFunc {
			return adminMiddleware
		},
	}
}
