package view

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/service"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

type UserRoute struct {
	service *service.UserService
}

func (u *UserRoute) UserSignIn(c *gin.Context) {
	req := request.ReqSignin{}
	req.Code = c.Query("code")
	req.State = c.Query("state")
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	data, err := tools.Base64Decode(req.State)
	if err != nil {
		c.HTML(http.StatusOK, config.GlobalConfig.Server.StaticPath+"/result.html", gin.H{"result": "签到失败:" + err.Error()})
		return
	}
	temp := service.MeetingMsg{}
	err = json.Unmarshal(data, &temp)
	if err != nil {
		c.HTML(http.StatusOK, config.GlobalConfig.Server.StaticPath+"/result.html", gin.H{"result": "签到失败:" + err.Error()})
		return
	}
	//validity test
	url, err := service.DefaultWsService.GetMeetingUrl(temp.MeetingID)
	if err != nil {
		c.HTML(http.StatusOK, config.GlobalConfig.Server.StaticPath+"/result.html", gin.H{"result": "签到失败:" + err.Error()})
		return
	}
	if url != temp.Code {
		c.HTML(http.StatusOK, config.GlobalConfig.Server.StaticPath+"/result.html", gin.H{"result": "签到失败: 二维码失效"})
		return
	}
	msg := service.SignCode{
		Code:      req.Code,
		Meeting:   temp.MeetingID,
		RetryTime: 0,
	}

	select {
	case u.service.SignMessage <- msg:
		c.HTML(http.StatusOK, config.GlobalConfig.Server.StaticPath+"/result.html", gin.H{"result": "签到成功"})
	default:
		c.HTML(http.StatusOK, config.GlobalConfig.Server.StaticPath+"/result.html", gin.H{"result": "签到失败,触发限流"})
	}
	return
}

func NewUserRoute() *UserRoute {
	route := &UserRoute{
		service: service.NewUserService(),
	}
	return route
}
