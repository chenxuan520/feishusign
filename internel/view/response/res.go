package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Err  string      `json:"err,omitempty"`
}

func Success(g *gin.Context, data interface{}) {
	g.JSON(http.StatusOK, Response{
		Code: 200,
		Data: data,
	})
}

func ResultHTML(g *gin.Context, status string, code int) {
	var resCode int
	switch code {
	case 0:
		resCode = http.StatusOK
	case 1:
		resCode = http.StatusServiceUnavailable
	case 2:
		resCode = http.StatusBadRequest
	}
	g.HTML(resCode, "result.html", map[string]interface{}{
		"code":   code,
		"status": status,
		"time":   time.Now().Format("1.2 15:04:05"),
	})
}

func Error(g *gin.Context, status int, err error) {
	g.JSON(status, Response{
		Code: status,
		Data: nil,
		Err:  err.Error(),
	})
}

func ErrorHTML(g *gin.Context, status int, err error) {
	g.HTML(status, "error.html", map[string]interface{}{
		"code": status,
		"err":  err.Error(),
	})
}
