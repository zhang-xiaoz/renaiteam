package main

import (
	"web/config"
	"web/log"
	"web/models"
	"web/router"
)

func main() {
	config.InitConfig() //初始化配置文件
	log.InitLogger()    //初始化日志文件
	log.SugarLogger.Info("配置文件读取完毕")
	Router := router.Routers() //分组路由器
	err := Router.Run(models.Overall_Situation_Server.Server.Port)
	if err != nil {
		log.SugarLogger.Error("启动服务器,端口%s失败", models.Overall_Situation_Server.Server.Port)
		log.SugarLogger.Panic("错误原因:", err)
	}
}
