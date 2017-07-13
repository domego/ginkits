package middleware

import "github.com/gin-gonic/gin"

var (
	AccessControlAllowOrigin = "*"
)

func CommonHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(AccessControlAllowOrigin) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", AccessControlAllowOrigin)
		}
		c.Next()
	}
}
