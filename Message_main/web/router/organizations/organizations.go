package organizations

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func Organizations(Router *gin.RouterGroup) { //社团信息
	{ //成员信息
		Router.POST("/get_member_message", api.Get_Member_Message) //?p=1&pn=6 //获取成员部分数据根据年级分类
	}
	{ //主攻方向
		Router.POST("/get_club_direction_message", api.Get_Club_Direction_Message) //获取主攻方向
	}
	{ //获奖信息
		Router.POST("/get_prize_message", api.Get_Prize_Message) //?p=1&pn=6 //获取获奖信息
	}
}
