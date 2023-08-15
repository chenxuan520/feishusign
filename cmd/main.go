package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view"
	"net/http"
	"os"
	"time"
)

func main() {
	initTime()
	initTest()
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
	err = logger.InitLog("Debug", "console", "[dian-feishu]", "logs", false, "LowercaseLevelEncoder", true)
	if err != nil {
		panic(err)
	}
	err = tools.InitLarkClient(config.GlobalConfig.Feishu.AppID, config.GlobalConfig.Feishu.AppSecret)
	if err != nil {
		panic(err)
	}
	err = model.InitMysql(config.GlobalConfig.Mysql.Dns())
	if err != nil {
		panic(err)
	}

	g := gin.Default()

	g.LoadHTMLGlob(config.GlobalConfig.Server.StaticPath + "/*")
	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://passport.feishu.cn/accounts"+
			"/auth_login/oauth2/authorize?client_id="+config.GlobalConfig.Feishu.AppID+
			"&redirect_uri=http://203.34.152.3:5204/api/admin/login&response_type=code&state=state123456")
	})

	g.StaticFile("/index", config.GlobalConfig.Server.StaticPath+"/index.html")

	view.InitGin(g)

	g.Run(fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))
}

// 时间修正
func initTime() {
	local := time.FixedZone("UTC +8:00", 8*3600)
	time.Local = local
}

func initTest() {
	if len(os.Args) < 2 {
		config.TestMode = false
	} else {
		if os.Args[1] == "-t" {
			config.TestMode = true
		}
	}
	return
}
