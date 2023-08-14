package view

import (
	"encoding/json"
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/service"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

type UserRoute struct {
	service *service.UserService
}

func (u *UserRoute) UserSignIn(c *gin.Context) {
	var msg service.SignCode
	req := request.ReqSignin{}
	req.Code = c.Query("code")
	req.State = c.Query("state")
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	data, err := tools.Base64Decode(req.State)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	temp := service.MeetingMsg{}
	err = json.Unmarshal(data, &temp)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	//validity test
	url, err := service.DefaultWsService.GetMeetingUrl(temp.MeetingID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	if url != temp.Code {
		response.ResultHTML(c, "签到失败,二维码失效", 1)
		return
	}
	msg = service.SignCode{
		Code:      req.Code,
		Meeting:   temp.MeetingID,
		RetryTime: 0,
	}
	select {
	case u.service.SignMessage <- msg:
		//response.Success(c, "success")
		response.ResultHTML(c, "签到成功", 0)
	default:
		response.ResultHTML(c, "签到失败,触发限流", 1)
	}
	return
}

func NewUserRoute() *UserRoute {
	route := &UserRoute{
		service: service.NewUserService(),
	}
	return route
}
