package Homeinformation

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func Homeinformation(Router *gin.RouterGroup) {
	Router.GET("/get/blog/picture:filename", api.Get_Blog_Picture) //获取博客首页的三张图片

}
