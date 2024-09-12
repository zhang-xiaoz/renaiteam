package router

import (
	"web/middlewares"
	"web/router/logon"
	"web/router/unlogon"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())    //跨域
	router1 := Router.Group("/logon") //登陆后功能
	{
		//router1.Use(middlewares.JWT_Check())
		logon.Logon(router1)
	}
	router2 := Router.Group("/unlogon") //未登录功能
	{
		unlogon.UnLogon(router2)
	}
	return Router
}
