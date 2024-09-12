package recruitment

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func Recruitment(Router *gin.RouterGroup) { //社团信息
	{ //培训时间加载
		Router.POST("/get_training_time", api.Get_Training_Time)
	}
	{ //培训人员
		Router.POST("/get_training_personnel", api.Get_Training_Personnel) //?p=1&pn=6
	}
	{ //培训信息
		Router.POST("/get_training_message", api.Get_Training_Message)
	}
}
