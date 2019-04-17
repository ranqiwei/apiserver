package middleware

import (
	"github.com/gin-gonic/gin"
	"apiserver/pkg/token"
	"apiserver/handler"
	"apiserver/pkg/errno"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
