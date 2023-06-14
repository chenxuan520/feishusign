package view

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/service"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/request"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view/response"
)

type AdminRoute struct {
	wsUpGrader *websocket.Upgrader
	service    *service.AdminService
}

func (a *AdminRoute) AdminLogin(c *gin.Context) {
	req := request.ReqSignin{}
	req.Code = c.Query("code")
	req.State = c.Query("state")
	if req.Code == "" || req.State == "" {
		response.Error(c, http.StatusBadRequest, fmt.Errorf("please log first"))
		return
	}
	jwt, err := a.service.AdminLogin(req.Code)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, map[string]interface{}{
		"jwt": jwt,
	})
}

func (a *AdminRoute) GetMeetingUrl(c *gin.Context) {
	//TODO jwt check

	resHeader := http.Header{}
	resHeader.Set("Sec-Websocket-Protocol", c.Request.Header.Get("Sec-Websocket-Protocol"))
	wsConn, err := a.wsUpGrader.Upgrade(c.Writer, c.Request, resHeader)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	ws := service.NewWsService(wsConn)
	go ws.Serve()
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
		service: service.NewAdminService(),
	}
	return &admin
}
