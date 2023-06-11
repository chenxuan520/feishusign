package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/larksuite/oapi-sdk-go/v3"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
)

var client *lark.Client

type ReqSignin struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

func TestSignIn(c *gin.Context) {
	req := ReqSignin{}
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"err": err.Error(),
		})
		return
	}
	fmt.Println(req)
	c.JSON(http.StatusOK, map[string]interface{}{
		"data": req,
	})
}

func main() {
	g := gin.Default()
	g.Use(gin.Recovery())
	client = lark.NewClient(config.GlobalConfig.Feishu.AppID, config.GlobalConfig.Feishu.AppSecret)
	api := g.Group("/api")
	api.GET("/ping", func(c *gin.Context) { c.JSON(http.StatusOK, map[string]interface{}{"ping": "pong"}) })
	api.GET("/user/signin", TestSignIn)
	g.Run(":5204")
}
