package kits

import (
	"github.com/gin-gonic/gin"

	"github.com/domego/ginkits/errors"
	//	"github.com/domego/gokits/log"
)

const (
	// ErrorCodeArgumentLack 缺少参数
	ErrorCodeArgumentLack = 100
	// ErrorCodeArgumentTypeInvalid 参数类型错误
	ErrorCodeArgumentTypeInvalid = 101
	// ErrorCodeInvalidSignature 签名错误
	ErrorCodeInvalidSignature = 102
	// ErrorCodeExpiredRequest 请求已过期
	ErrorCodeExpiredRequest = 103
)

type RespErrorMessage struct {
	Code    int    `json:"code"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Message string `json:"message"`
}

// RenderSuccess 响应成功的数据
func RenderSuccess(c *gin.Context, data interface{}) {
	resp := gin.H{
		"success": true,
		"result":  data,
	}
	// log.Tracef("%+v", resp)
	c.JSON(200, resp)
}

// RenderError 响应错误信息
func RenderError(c *gin.Context, err interface{}) {
	resp := gin.H{
		"success": false,
		"error":   err,
	}
	// log.Infof("%+v", resp)
	c.JSON(200, resp)
	c.Abort()
}

func RenderErrorMessage(c *gin.Context, code string) {
	err := errorkits.Get(code)
	if err == nil {
		err = &errorkits.ErrorMessage{
			Content: "网络不给力",
		}
	}
	RenderError(c, err)
}
