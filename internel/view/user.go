package view

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

func UserSignIn(c *gin.Context) {
	req := request.ReqSignin{}
	err := c.BindQuery(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	fmt.Println(req)
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": req,
	})
}
