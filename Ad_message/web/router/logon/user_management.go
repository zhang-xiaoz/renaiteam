package logon

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func User_Management(Router *gin.RouterGroup) {
	router1 := Router.Group("/register")
	{
		router1.POST("/agree_register_user", api.Agree_Register_User)   //用户注册申请通过
		router1.POST("/get_register_user", api.Get_Register_User)       //?p=1&pn=6 //注册的用户获取//分页查找
		router1.POST("/refuse_register_user", api.Refuse_Register_User) //用户注册申请拒绝//需要发邮箱
		router1.POST("/add_blacklist", api.Add_Register_Blacklist)      //拉入黑名单
	}
	router2 := Router.Group("/revise")
	{
		router2.POST("/get_user", api.Get_User)               //?p=1&pn=6 //所有用户信息获取//获取部分可以查看信息
		router2.POST("/revise_user", api.Revise_User)         //修改用户数据
		router2.POST("/delete_user", api.Delete_Revise_User)  //删除用户数据//删除这个用户所有数据
		router2.POST("/seek_grade_user", api.Seek_Grade_User) //?p=1&pn=6 //根据年级查找
		router2.POST("/seek_name_user", api.Seek_Name_User)   //?p=1&pn=6//根据姓名查找
	}
	router3 := Router.Group("/cancel")
	{
		router3.POST("/get_user", api.Get_Cancel_User)       //?p=1&pn=6获取注销用户信息
		router3.POST("/delete_user", api.Delete_Cancel_User) //删除注销用户
	}
	router4 := Router.Group("/blacklist") //黑名单
	{
		router4.POST("/get_user", api.Get_Blacklist_User)       //?p=1&pn=6获取黑名单邮箱
		router4.POST("/delete_user", api.Delete_Blacklist_User) //移除黑名单
		router4.POST("/add_user", api.Add_Blacklist_User)       //增加黑名单
	}
}
