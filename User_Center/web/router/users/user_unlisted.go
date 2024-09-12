package users

import (
	"web/api"
	"web/middlewares"

	"github.com/gin-gonic/gin"
)

func User_Unlisted(Router *gin.RouterGroup) {
	//注册时填写邮箱，姓名，用户名(昵称)，密码，邮箱验证码(一分钟只能按一次且五分钟之内可以使用)
	{
		//未登录时注册时邮箱验证码获取
		Router.POST("/get/mailbox/code", api.User_Unlisted_Get_Mailbox_Code)
		//未登录时注册时进行注册//进行验证
		Router.POST("/register/mailbox", api.User_Unlisted_Register_Mailbox)
		//未登录时进行登录
		Router.POST("/login", middlewares.Picture_Check(), api.User_Unlisted_Login)
		//获取图片验证
		Router.POST("/get/picture/code", api.User_Unlisted_Get_Picture_Code)
		//未注册时获取验证码(用于更改密码)
		Router.POST("/get/code", api.User_Unlisted_Get_Code)
		//未注册时获取验证码进行验证(用于更改密码)
		Router.POST("/check/code", api.User_Unlisted_Check_Code)
		//未登录时密码更改(验证码已经通过验证)
		Router.POST("/change/password", api.User_Unlisted_Change_Password)
	}
}
