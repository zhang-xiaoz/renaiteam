package logon

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func Front_End_Data(Router *gin.RouterGroup) {
	router1 := Router.Group("/home_data") //首页三数据管理
	{
		router1.POST("/revise_message_mysql", api.Revise_Message_Mysql) //关于我们数据修改（文字和社团位置修改（文字加图片）学习方式
	}
	router2 := Router.Group("/club_recruitment") //社团招新数据管理
	{
		router2.POST("/revise_training_time", api.Revise_Training_Time)    //培训时间修改//时间格式2024-06-03
		router2.POST("/add_training_time", api.Add_Training_Time)          //培训时间增加//先查询后增加
		router2.POST("/delete_training_time", api.Delete_Training_Time)    //培训时间删除//先删除完后从新增加
		router1.POST("/revise_training_message", api.Revise_Message_Mysql) //培训信息修改
	}
	router3 := Router.Group("/community_information") //社团信息数据管理
	{
		router3.POST("/revise_club_direction", api.Revise_Message_Mysql)        //修改主攻方向
		router3.POST("/delete_club_direction", api.Delete_Club_Direction)       //删除主攻方向
		router3.POST("/add_club_direction", api.Add_Club_Direction)             //增加主攻方向
		router3.POST("/revise_award_information", api.Revise_Award_Information) //获奖信息修改
		router3.POST("/delete_award_information", api.Delete_Award_Information) //获奖信息删除
		router3.POST("/add_award_information", api.Add_Award_Information)       //获奖信息增加
	}
}
