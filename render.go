package kits

import "github.com/gin-gonic/gin"

// RenderSuccess 响应成功的数据
func RenderSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

// RenderError 响应错误信息
func RenderError(c *gin.Context, err interface{}) {
	c.JSON(200, map[string]interface{}{
		"success": false,
		"error":   err,
	})
}
