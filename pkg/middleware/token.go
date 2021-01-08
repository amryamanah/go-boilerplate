package middleware

import (
	"github.com/amryamanah/go-boilerplate/internal/auth"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Printf("[TokenAuthMiddleware] Request: %+v", ctx.Request)
		err := auth.TokenValid(ctx.Request)
		if err != nil {
			log.Println("[TokenAuthMiddleware] error")
			ctx.JSON(http.StatusUnauthorized, err.Error())
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
