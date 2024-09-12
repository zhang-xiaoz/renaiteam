package logon

import "github.com/gin-gonic/gin"

func Logon(Router *gin.RouterGroup) { //登录后功能
	router1 := Router.Group("/front_end_data") //前端页面数据管理
	{
		Front_End_Data(router1)
	}
	router2 := Router.Group("/user_management") //用户管理
	{
		User_Management(router2)
	}
	//博客信息管理
	router3 := Router.Group("/blog_information") //用户管理
	{
		Blog_Information(router3)
	}
}
