package params

import (
	kits "github.com/domego/ginkits"
	"github.com/domego/gokits"
	"github.com/gin-gonic/gin"
)

func Required(c *gin.Context, name string, value interface{}) bool {
	if utils.IsEmpty(value) {
		kits.RenderError(c, &kits.RespErrorMessage{
			Code:    kits.ErrorCodeArgumentLack,
			Message: name + " is required",
		})
		return false
	}
	return true
}
