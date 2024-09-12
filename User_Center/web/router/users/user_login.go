package users

import (
	"web/api"
	"web/middlewares"

	"github.com/gin-gonic/gin"
)

func User_Login(Router *gin.RouterGroup) {
	//登陆后个人信息
	{
		Router.POST("/get/user/message", middlewares.JWT_Check(), api.User_Login_Get_User_Message)                     //获取基本信息（暂未添加年级，方向，qq,微信，职位，座右铭）
		Router.POST("/change/user/message", middlewares.JWT_Check(), api.User_Login_Change_User_Message)               //修改基本信息
		Router.POST("/change/user/mailbox/code", middlewares.JWT_Check(), api.User_Login_Change_User_Mailbox_Code)     //修改邮箱(验证现在邮箱)
		Router.POST("/get/user/mailbox/code", middlewares.JWT_Check(), api.User_Login_Get_User_Malibox_Code)           //真正修改邮箱是获取验证码
		Router.POST("/change/user/mailbox/revise", middlewares.JWT_Check(), api.User_Login_Change_User_Mailbox_Revise) //修改邮箱(真正修改)
		Router.POST("/change/user/code", middlewares.JWT_Check(), api.User_Login_Change_User_Code)                     //获取验证码
		Router.POST("/change/user/password", middlewares.JWT_Check(), api.User_Login_Change_User_Password)             //修改密码(验证码,加锁(登录时的锁))
		Router.POST("/logout/user", middlewares.JWT_Check(), api.User_Login_Logout_User)                               //注销账户(验证码，注销理由)
		Router.POST("/quit", api.User_Quit)                                                                            //退出登录
	}
	{ //头像处理
		Router.POST("/get/user/picture/name", middlewares.JWT_Check_New(), api.Get_User_Picture_Name)        //获取头像名字
		Router.GET("/get/user/picture/:filename", api.User_Login_Get_User_Picture)                           //获取头像
		Router.POST("/change/user/picture", middlewares.JWT_Check_New(), api.User_Login_Change_User_Picture) //修改头像
	}
}

//现在我将用户注销信息存储到mysql里边的password里边
