package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Err  string      `json:"err,omitempty"`
}

// Success 成功
func Success(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
	})
}

// Error 错误
func Error(g *gin.Context, status int, err error) {
	g.JSON(status, Response{
		Code: status,
		Data: nil,
		Err:  err.Error(),
	})
}

// Error 错误
func ErrorDetail(g *gin.Context, status int, data interface{}, err error) {
	res := Response{
		Code: status,
		Data: data,
		Err:  err.Error(),
	}

	g.JSON(status, res)
}
