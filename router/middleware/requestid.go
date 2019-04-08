package middleware

import "github.com/gin-gonic/gin"
import "github.com/satori/go.uuid"

//global middleware
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Check for incoming header,use it if exists
		requestId := c.Request.Header.Get("X-Request-Id")
		//create uuid with uuid4
		if requestId == "" {
			u4, _ := uuid.NewV4()
			requestId = u4.String()
		}

		//expose it for use in the application
		c.Set("X-Request-Id", requestId)
		//set X-Request-Id header to response
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}
