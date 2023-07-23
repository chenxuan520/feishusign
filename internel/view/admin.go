package view

import (
	"fmt"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
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
	//response.Success(c, map[string]interface{}{
	//	"jwt": jwt,
	//})

	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"jwt": jwt,
	})
}

func (a *AdminRoute) GetMeetingUrl(c *gin.Context) {
	// jwt check and get user id
	token := c.Query("jwt")
	if token == "" {
		response.ErrorHTML(c, http.StatusBadRequest, fmt.Errorf("no jwt token checked"))
		return
	}
	claims, err := tools.ParseJwtToken(token)
	if err != nil {
		response.ErrorHTML(c, http.StatusBadRequest, err)
		return
	}
	userid := claims.UserId

	meeting := c.Query("meeting")
	if meeting == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please take meeting query"))
		return
	}
	// check if exist meeting
	_, err = model.GetMeetingByID(meeting)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("meeting don't exist"))
		return
	}

	//upgrade to websocket
	resHeader := http.Header{}
	resHeader.Set("Sec-Websocket-Protocol", c.Request.Header.Get("Sec-Websocket-Protocol"))
	wsConn, err := a.wsUpGrader.Upgrade(c.Writer, c.Request, resHeader)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	//add conn
	err = a.WsService.AddWsConn(wsConn, userid, meeting)
	if err != nil {
		return
	}
}

func (a *AdminRoute) CreateMeeting(c *gin.Context) {
	// check jwt token
	token := c.Query("jwt")
	if token == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("no jwt token checked"))
		return
	}
	claims, err := tools.ParseJwtToken(token)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	userId := claims.UserId

	// create meeting
	str, err := a.adminService.AdminCreateMeeting(userId)
	if err != nil {
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
