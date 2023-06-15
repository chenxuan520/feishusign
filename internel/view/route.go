package view

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-gin"
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
			meetingGroup.GET("/url", adminRoute.GetMeetingUrl)
			//TODO
			meetingGroup.GET("/latest")
		}
	}

	//event
	//TODO finish it
	eventRoute := NewEventRoute()
	api.POST("/event", sdkginext.NewEventHandlerFunc(eventRoute.InitEvent()))
}

func initMiddle(g *gin.Engine) {
	g.Use(gin.Recovery())
}
