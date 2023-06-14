package view

import (
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
		//meeting
		meeting := adminGroup.Group("/meeting")
		{
			//TODO
			meeting.GET("/url", adminRoute.GetMeetingUrl)
		}
	}

}

func initMiddle(g *gin.Engine) {
	g.Use(gin.Recovery())
}
