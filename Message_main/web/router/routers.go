package router

import (
	"web/middlewares"
	"web/router/Homeinformation"
	"web/router/homepage"
	"web/router/organizations"
	"web/router/recruitment"

	"github.com/gin-gonic/gin"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middlewares.Cors())              //跨域
	router1 := Router.Group("/Homeinformation") //博客首页图片信息
	{
		Homeinformation.Homeinformation(router1)
	}
	router2 := Router.Group("/organizations_mes") //社团信息
	{
		organizations.Organizations(router2)
	}
	router3 := Router.Group("/club_recruitment") //社团招新信息
	{
		recruitment.Recruitment(router3)
	}
	router4 := Router.Group("/homepage") //首页三个数据加载
	{
		homepage.Homepage(router4)
	}
	return Router
}
