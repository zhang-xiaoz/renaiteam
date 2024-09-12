package main

import (
	"web/config"
	"web/log"
	"web/models"
	"web/router"
	"web/tool"
)

func main() {
	config.InitConfig()        //初始化配置文件
	log.InitLogger()           //初始化日志文件//日志文件按日分开//日志文件保留一个星期
	tool.InitVaildators()      //注册字符验证器
	Router := router.Routers() //分组路由器
	err := Router.Run(models.Overall_Situation_Server.Server.Port)
	if err != nil {
		log.SugarLogger.Error("启动服务器,端口%s失败", models.Overall_Situation_Server.Server.Port)
		log.SugarLogger.Panic("错误原因:", err)
	}
}
