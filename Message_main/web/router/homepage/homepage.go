package homepage

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func Homepage(Router *gin.RouterGroup) { //首页信息加载
	{ //关于我们
		Router.POST("/get_about_us", api.Get_About_Us)
	}
	{ //学习方式
		Router.POST("/get_learning_style", api.Get_Learning_Style)
	}
	{ //社团位置
		Router.POST("/get_club_location", api.Get_Club_Location)
	}
}
