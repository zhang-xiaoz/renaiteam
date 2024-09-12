package users

import (
	"github.com/gin-gonic/gin"
)

func User(Router *gin.RouterGroup) {
	router1 := Router.Group("/unlisted") //未登录
	{
		User_Unlisted(router1)
	}
	router2 := Router.Group("/listed") //登录后
	{
		User_Login(router2)
	}
}
