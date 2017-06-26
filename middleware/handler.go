package middleware

import "github.com/gin-gonic/gin"

func CommonHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
