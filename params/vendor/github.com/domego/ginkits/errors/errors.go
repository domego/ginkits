package errorkits

import (
	"io/ioutil"
	"os"

	"github.com/domego/gokits/log"
	"github.com/kr/pretty"
	yaml "gopkg.in/yaml.v1"
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
	// ErrorCodeArgumentValueInvalid 非法参数值
	ErrorCodeArgumentValueInvalid = 104
)

type ErrorMessage struct {
	Code    int    `json:"code" yaml:"code"`
	Title   string `json:"title" yaml:"title"`
	Content string `json:"content" yaml:"content"`
	Message string `json:"message" yaml:"message"`
}

var (
	ErrorMessageFile = "config/error_msg.yaml"
	errorMessageMap  = map[string]*ErrorMessage{}
)

func ReadConfig(exit bool) map[string]*ErrorMessage {
	bs, err := ioutil.ReadFile(ErrorMessageFile)
	if err != nil {
		log.Errorf("load config file failed, the path is %s", ErrorMessageFile)
		if exit {
			os.Exit(1)
		}
	}
	config := map[string]*ErrorMessage{}
	err = yaml.Unmarshal(bs, &config)
	if err != nil {
		log.Errorf("unmarshal config failed, %s", err)
		if exit {
			os.Exit(1)
		}
	}
	log.Debugf("read config: %+v", pretty.Formatter(config))
	return config
}

func Init() map[string]*ErrorMessage {
	errorMessageMap = ReadConfig(true)
	return errorMessageMap
}

func Get(key string) *ErrorMessage {
	return errorMessageMap[key]
}
