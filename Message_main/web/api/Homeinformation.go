package api

import (
	"os"

	"github.com/gin-gonic/gin"
)

func Get_Blog_Picture(ctx *gin.Context) {
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	// 获取文件名
	filename := ctx.Param("filename")
	// 图片所在的文件夹路径
	dir := "./img/blog/"

	// 完整的图片路径
	filePath := dir + filename

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
			"jwt":     jwtString,
		})
		return
	}
	// 使用Gin的c.File方法发送文件
	ctx.File(filePath)
}
