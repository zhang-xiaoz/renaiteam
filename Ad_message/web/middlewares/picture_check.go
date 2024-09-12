package middlewares

import (
	"web/log"
	"web/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 验证邮箱验证码是否正确
func Picture_Check() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//实体化对象
		login := models.User_login{}
		err := ctx.ShouldBindBodyWith(&login, binding.JSON) //这是一个通用方法
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "验证码不正确",
			})
			ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
			return
		}
		m := models.Store.Verify(login.Id, login.B64s, true)
		if !m {
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "验证码不正确",
			})
			ctx.Abort()
			return
		}
		ctx.Next() //直接继续往下走
	}
}
