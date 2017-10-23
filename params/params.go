package paramkits

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

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
	data, _ := ioutil.ReadAll(c.Request.Body)
	jsonData := gjson.Parse(string(data))
	if jsonData.IsObject() {
		jsonData.ForEach(func(key, value gjson.Result) bool {
			params[key.String()] = value.String()
			return true
		})
	}
	if c.Request.PostForm == nil || c.Request.Form == nil {
		c.Request.ParseForm()
	}
	for k, v := range c.Request.Form {
		params[k] = v[0]
	}
	return params
}

// ClientInfo 用户快速获取访问信息
type ClientInfo struct {
	UserAgent string
	IP        string
}

// ParseClientInfo new client info
func ParseClientInfo(c *gin.Context) ClientInfo {
	return ClientInfo{
		UserAgent: c.Request.UserAgent(),
		IP:        c.ClientIP(),
	}
}
