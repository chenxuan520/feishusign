package view

import (
	"fmt"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	"log"
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
			meetingGroup.GET("/url", adminRoute.GetMeetingUrl)
			meetingGroup.GET("/create", adminRoute.CreateMeeting)
		}
	}

	//event
	//TODO finish it

	//api.POST("/event", usedForEventDebug)
	eventRoute := NewEventRoute()
	api.POST("/event", sdkginext.NewEventHandlerFunc(eventRoute.InitEvent()))
}

func initMiddle(g *gin.Engine) {
	g.Use(gin.Recovery())
}

func usedForEventDebug(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		log.Println("here err :", err)
		return
	}
	fmt.Println(string(data))
	return
}
