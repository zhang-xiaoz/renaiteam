package logon

import (
	"web/api"

	"github.com/gin-gonic/gin"
)

func Blog_Information(Router *gin.RouterGroup) {
	//用户博客管理
	router1 := Router.Group("/blog")
	{
		router1.POST("/get_blog_name", api.Get_Blog_Name)   //?p=1&pn=6根据姓名查找
		router1.POST("/get_blog_title", api.Get_Blog_Title) //?p=1&pn=6根据标题查找
		router1.POST("/get_blog_label", api.Get_Blog_Label) //?p=1&pn=6根据标签查找
		router1.POST("/get_blog_all", api.Get_Blog_All)     // ?p=1&pn=6//博客首页查看
		router1.POST("/get_blog_one", api.Get_Blog_One)     //博客点开内容获取
		router1.POST("/delete_blog", api.Delete_Blog)       //删除博客内容（带上原因）//修改elastic
	}
	//博客审核管理
	router2 := Router.Group("/to_examine")
	{
		router2.POST("/get_article_review", api.Get_Article_Review)       // ?p=1&pn=6//获取审核的内容
		router2.POST("/one_article_review", api.One_Article_Review)       //获取单个内容
		router2.POST("/agree_article_review", api.Agree_Article_Review)   //同意审核
		router2.POST("/refuse_article_review", api.Refuse_Article_Review) //拒绝审核(理由)
	}
}
