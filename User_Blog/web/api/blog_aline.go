package api

import (
	"context"
	"web/log"
	"web/models"
	modelss "web/models/blogs"
	"web/proto_blog"
	"web/proto_web"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Quick_Response_State(ctx *gin.Context) { //快速判断响应状态
	//先获取请求的两个数据
	access_token := ctx.GetHeader("access_token")
	refresh_token := ctx.GetHeader("refresh_token")
	uuid := ctx.GetHeader("uuid")
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Host + models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto_web.NewUsersClient(conn)
	if access_token != "" {
		result, _ := userWebClient.Get_Redis_Storage(context.Background(), &proto_web.Redis_Storage{
			Key: uuid + "access_token",
			Who: "1",
		})
		if result != nil {
			if result.Value == access_token {
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "登录成功",
				})
				return
			}
		}
	}
	if refresh_token != "" {
		result, _ := userWebClient.Get_Redis_Storage(context.Background(), &proto_web.Redis_Storage{
			Key: uuid + "refresh_token",
			Who: "1",
		})
		if result != nil {
			if result.Value == refresh_token {
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "登录成功",
				})
				return
			}
		}
	}
	ctx.JSON(200, map[string]string{
		"code":    "200",
		"message": "登录过期",
	})
}

func Add_Collection_Blog(ctx *gin.Context) { //收藏博客
	//先获取请求的两个数据
	access_token := ctx.GetHeader("access_token")
	uuid := ctx.GetHeader("uuid")
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Host + models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto_web.NewUsersClient(conn)
	if access_token != "" {
		result, _ := userWebClient.Get_Redis_Storage(context.Background(), &proto_web.Redis_Storage{
			Key: uuid + "access_token",
			Who: "1",
		})
		if result != nil {
			if result.Value != access_token {
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录过期",
				})
				return
			}
		}
	} else {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	//开始获取文章uuid
	message := modelss.Blog_Delete{}
	//解析数据
	err = ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败1",
		})
		return
	}
	//开始存储到数据库里边
	//连接Grpc服务
	conn1, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败2",
		})
		return
	}
	defer conn1.Close()
	userblogClient := proto_blog.NewBlogsClient(conn1)
	_, err = userblogClient.Add_Mysql_Blog_Collection(context.Background(), &proto_blog.Blog_UserUuid{
		UserUuid: uuid,
		Status:   message.Article_Uuid,
	})
	if err != nil {
		//保存失败
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "保存失败3",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "保存成功",
	})
}

func Delete_Collection_Blog(ctx *gin.Context) { //删除收藏的博客
	//先获取请求的两个数据
	access_token := ctx.GetHeader("access_token")
	uuid := ctx.GetHeader("uuid")
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Host + models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto_web.NewUsersClient(conn)
	if access_token != "" {
		result, _ := userWebClient.Get_Redis_Storage(context.Background(), &proto_web.Redis_Storage{
			Key: uuid + "access_token",
			Who: "1",
		})
		if result != nil {
			if result.Value != access_token {
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录过期",
				})
				return
			}
		}
	} else {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	//开始获取文章uuid
	message := modelss.Blog_Delete{}
	//解析数据
	err = ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
		})
		return
	}
	//连接Grpc服务
	conn1, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
		})
		return
	}
	defer conn1.Close()
	//开始删除elastic里边的数据
	userblogClient := proto_blog.NewBlogsClient(conn1)
	_, err = userblogClient.Delete_Mysql_Blog_Collection(context.Background(), &proto_blog.Blog_UserUuid{
		UserUuid: uuid,
		Status:   message.Article_Uuid,
	})
	if err != nil {
		//保存失败
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "删除失败",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "删除成功",
	})
}

func Get_Collection_Status(ctx *gin.Context) { //获取收藏状态
	//先获取请求的两个数据
	access_token := ctx.GetHeader("access_token")
	uuid := ctx.GetHeader("uuid")
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Host + models.Overall_Situation_Grpc_Server_web.Webgrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto_web.NewUsersClient(conn)
	if access_token != "" {
		result, _ := userWebClient.Get_Redis_Storage(context.Background(), &proto_web.Redis_Storage{
			Key: uuid + "access_token",
			Who: "1",
		})
		if result != nil {
			if result.Value != access_token {
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录过期",
				})
				return
			}
		}
	} else {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录过期",
		})
		return
	}
	//开始获取文章id和用户id判断状态
	message := modelss.Blog_Delete{}
	//解析数据
	err = ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "获取失败",
		})
		return
	}
	//连接Grpc服务
	conn1, err := grpc.Dial((models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Host + models.Overall_Situation_Grpc_Server_blog.Bloggrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "获取失败",
		})
		return
	}
	defer conn1.Close()
	//开始判断状态
	//状态包括  文章是这个用户的  文章已经被这个用户收藏  文章没有被收藏
	userblogClient := proto_blog.NewBlogsClient(conn1)
	result1, err := userblogClient.Get_Blog_Collection_Status(context.Background(), &proto_blog.Blog_UserUuid{
		UserUuid: uuid,
		Status:   message.Article_Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		//获取失败
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "获取失败",
		})
		return
	}
	//返回状态
	if result1.Status == 1 {
		//获取失败
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "该用户是作者",
		})
	} else if result1.Status == 2 {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "该用户已经收藏",
		})
	} else if result1.Status == 3 {
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "该用户未收藏",
		})
	} else {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "获取失败",
		})
	}

}
