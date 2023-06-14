package view

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/service"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

type UserRoute struct {
	service *service.UserService
}

func (u *UserRoute) UserSignIn(c *gin.Context) {
	req := request.ReqSignin{
		Code:  "",
		State: "",
	}
	req.Code = c.Query("code")
	req.State = c.Query("state")
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	msg := service.SignCode{
		Code:      req.Code,
		Meeting:   0,
		RetryTime: 0,
	}
	select {
	case u.service.SignMessage <- msg:
		response.Success(c, map[string]interface{}{})
	default:
		response.Error(c, http.StatusBadRequest, fmt.Errorf("触发限流,稍后再试"))
	}
}

func NewUserRoute() *UserRoute {
	route := &UserRoute{
		service: service.NewUserService(),
	}
	return route
}
