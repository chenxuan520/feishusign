package view

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
)

func InitGin(g *gin.Engine) {
	api := g.Group("/api")
	api.GET("ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, map[string]interface{}{"ping": "pong"}) })

	//user
	userRoute := NewUserRoute()
	user := api.Group("/user")
	user.GET("/signin", userRoute.UserSignIn)

	//admin
	admin := api.Group("/admin")
	meeting := admin.Group("/meeting")
	meeting.POST("/create")
	//TODO
	meeting.GET("/url", func(ctx *gin.Context) {
		str := url.QueryEscape(config.GlobalConfig.Server.RedirectURL)
		stdUrl := fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s&state=%s", str, config.GlobalConfig.Feishu.AppID, "hello")
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"url": stdUrl,
		})
	})
}

func initMiddle(g *gin.Engine) {
	g.Use(gin.Recovery())
}
