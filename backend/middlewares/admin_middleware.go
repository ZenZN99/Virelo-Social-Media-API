package middlewares

import (
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		user, exists := ctx.Get("user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			ctx.Abort()
			return
		}

		payload, ok := user.(*utils.TokenPayload)
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user data",
			})
			ctx.Abort()
			return
		}

		if payload.Role != "admin" {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error": "Only admins can access this operation",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
