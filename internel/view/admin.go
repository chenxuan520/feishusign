package view

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/middlerware"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/service"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

type AdminRoute struct {
	wsUpGrader   *websocket.Upgrader
	adminService *service.AdminService
	WsService    *service.WsService
}

func (a *AdminRoute) AdminLogin(c *gin.Context) {
	req := request.ReqSignin{}
	req.Code = c.Query("code")
	req.State = c.Query("state")
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	jwt, err := a.adminService.AdminLogin(req.Code)
	if err != nil {
		response.ErrorHTML(c, http.StatusBadRequest, err)
		return
	}

	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"jwt": jwt,
	})
}

func (a *AdminRoute) GetMeetingUrl(c *gin.Context) {
	// jwt check and get user id
	token := c.Query("jwt")
	if token == "" {
		response.Error(c, http.StatusUnauthorized, fmt.Errorf("no jwt token checked"))
		return
	}
	userId, err := middlerware.VerifyJwt(token)
	if err != nil {
		response.Error(c, http.StatusForbidden, err)
		return
	}

	meeting := c.Query("meeting")
	if meeting == "" {
		err := fmt.Errorf("no meeting query parameter found")
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	// check if exist meeting

	if _, err := model.GetMeetingByID(meeting); err != nil {
		if err == model.NotFind {
			err = fmt.Errorf("meeting don't exist")
		}
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	//upgrade to websocket
	resHeader := http.Header{}
	resHeader.Set("Sec-Websocket-Protocol", c.Request.Header.Get("Sec-Websocket-Protocol"))
	wsConn, err := a.wsUpGrader.Upgrade(c.Writer, c.Request, resHeader)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	//add conn
	err = a.WsService.AddWsConn(wsConn, userId, meeting)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		return
	}
}

func (a *AdminRoute) CreateMeeting(c *gin.Context) {
	rawUserId, exists := c.Get("uid")
	if !exists {
		err := fmt.Errorf("no userId found in context")
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	// 需要将type any转化为type string
	userId := fmt.Sprintf("%v", rawUserId)

	// create meeting
	str, err := a.adminService.AdminCreateMeeting(userId)
	if err != nil {
		logger.GetLogger().Error(err.Error())
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, map[string]interface{}{"meeting": str})
}

func NewAdminRoute() *AdminRoute {
	admin := AdminRoute{
		wsUpGrader: &websocket.Upgrader{
			HandshakeTimeout: 10 * time.Second,
			ReadBufferSize:   1024 * 4,
			WriteBufferSize:  1024 * 4,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		adminService: service.NewAdminService(),
		WsService:    service.NewWsService(),
	}
	return &admin
}
