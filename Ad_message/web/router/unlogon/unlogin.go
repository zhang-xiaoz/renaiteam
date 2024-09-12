package unlogon

import (
	"web/api"
	"web/middlewares"

	"github.com/gin-gonic/gin"
)

func UnLogon(Router *gin.RouterGroup) { //未登录功能
	Router.POST("/get_code_message", api.Get_Code_Message)                    //验证码获取
	Router.POST("/check_logon", middlewares.Picture_Check(), api.Check_Logon) //登录验证
}
