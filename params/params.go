package params

// BaseParam 请求基础参数
type BaseParam struct {
	OS      string `form:"os" json:"os"`
	Net     string `form:"net" json:"net"`
	Channel string `form:"channel" json:"channel"`
	AppVer  string `form:"app_ver" json:"app_ver"`
}
