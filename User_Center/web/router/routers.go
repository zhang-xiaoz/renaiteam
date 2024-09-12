package router

import (
	"web/middlewares"
	"web/router/ad"
	"web/router/users"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())   //跨域
	router1 := Router.Group("/user") //用户
	{
		users.User(router1)
	}
	router2 := Router.Group("/ad") //管理员
	{
		ad.Ad(router2)
	}
	return Router
}
