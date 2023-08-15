package view

import (
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/middlerware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitGin(g *gin.Engine) {
	api := g.Group("/api")
	api.GET("ping", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, map[string]interface{}{"ping": "pong"}) })

	//user
	userRoute := NewUserRoute()
	userGroup := api.Group("/user")
	{
		userGroup.GET("/signin", userRoute.UserSignIn)
	}

	//admin
	adminRoute := NewAdminRoute()
	adminGroup := api.Group("/admin")
	{
		//login
		adminGroup.GET("/login", adminRoute.AdminLogin)
		//meetingGroup
		meetingGroup := adminGroup.Group("/meeting")
		{
			// 用于测试其他功能，跳过校验
			// meetingGroup.Use(middlerware.Debug())

			meetingGroup.GET("/url", adminRoute.GetMeetingUrl)
			meetingGroup.GET("/create", middlerware.Auth(), adminRoute.CreateMeeting)
		}
	}

	//event
	eventRoute := NewEventRoute()
	api.POST("/event", sdkginext.NewEventHandlerFunc(eventRoute.InitEvent()))
}

func initMiddle(g *gin.Engine) {
	g.Use(gin.Recovery())
}
