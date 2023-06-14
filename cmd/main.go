package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/config"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/logger"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/model"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/tools"
	"gitlab.dian.org.cn/dianinternal/feishusign/internel/view"
)

func main() {
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
	view.InitGin(g)

	g.Run(fmt.Sprintf(":%d", config.GlobalConfig.Server.Port))
}
