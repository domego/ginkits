package paramkits

import "github.com/gin-gonic/gin"

// BaseParam 请求基础参数
type BaseParam struct {
	OS      string `form:"os" json:"os"`
	Net     string `form:"net" json:"net"`
	Channel string `form:"channel" json:"channel"`
	AppVer  string `form:"app_ver" json:"app_ver"`
}

// GetParams 获取所有请求参数
func GetParams(c *gin.Context) map[string]string {
	params := make(map[string]string)
	for k, v := range c.Request.Form {
		params[k] = v[0]
	}
	return params
}
