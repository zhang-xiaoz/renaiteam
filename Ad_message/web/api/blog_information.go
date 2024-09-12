package api

import (
	"context"
	"encoding/json"
	"strconv"
	"web/log"
	"web/models"
	"web/proto_blog"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Get_Blog_All(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Manage_All(context.Background(), &proto_blog.Blog_Page{
		P:       p,
		Pn:      pn,
		Message: "已经审核",
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	log.SugarLogger.Error(err)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result.Message,
	})
}

func Get_Blog_One(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Blog_Uuid{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Mysql_Message(context.Background(), &proto_blog.Blog_ArticleUuid{
		ArticleUuid: message1.Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	jsondata, _ := json.Marshal(result)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    string(jsondata),
	})
}

func Delete_Blog(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Blog_Delete{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	_, err = blogWebClient.Delete_El6_Mysql_Process(context.Background(), &proto_blog.BlogProcessBlog{
		ArticleUuid: message1.Uuid,
		Process:     "草稿",
		UserUuid:    message1.Reasion,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "删除成功",
		"jwt":     jwtString,
	})
}

func Get_Article_Review(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Manage_All(context.Background(), &proto_blog.Blog_Page{
		P:       p,
		Pn:      pn,
		Message: "未审核",
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	log.SugarLogger.Error(err)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result.Message,
	})
}

func One_Article_Review(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Blog_Uuid{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Mysql_Message(context.Background(), &proto_blog.Blog_ArticleUuid{
		ArticleUuid: message1.Uuid,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	jsondata, _ := json.Marshal(result)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    string(jsondata),
	})
}

func Agree_Article_Review(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Blog_Delete{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	_, err = blogWebClient.Delete_El6_Mysql_Process(context.Background(), &proto_blog.BlogProcessBlog{
		ArticleUuid: message1.Uuid,
		Process:     "已经审核",
		UserUuid:    "审核通过",
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "成功",
		"jwt":     jwtString,
	})

}

func Refuse_Article_Review(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Blog_Delete{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	_, err = blogWebClient.Delete_El6_Mysql_Process(context.Background(), &proto_blog.BlogProcessBlog{
		ArticleUuid: message1.Uuid,
		Process:     "草稿",
		UserUuid:    message1.Reasion,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "成功",
		"jwt":     jwtString,
	})
}

func Get_Blog_Name(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Message{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Manage_Name(context.Background(), &proto_blog.Blog_Page{
		P:       p,
		Pn:      pn,
		Message: message1.Message,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	jsondata, _ := json.Marshal(result)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    string(jsondata),
	})
}

func Get_Blog_Title(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Message{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Manage_Title(context.Background(), &proto_blog.Blog_Page{
		P:       p,
		Pn:      pn,
		Message: message1.Message,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	jsondata, _ := json.Marshal(result)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    string(jsondata),
	})
}

func Get_Blog_Label(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Message{}
	err := ctx.ShouldBindBodyWith(&message1, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接bloggrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Host + models.Overall_Situation_Blog_Grpc.BlogGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	defer conn.Close()
	blogWebClient := proto_blog.NewBlogsClient(conn)
	//开始获取数据
	result, err := blogWebClient.Get_Blog_Manage_Label(context.Background(), &proto_blog.Blog_Page{
		P:       p,
		Pn:      pn,
		Message: message1.Message,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	jsondata, _ := json.Marshal(result)
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    string(jsondata),
	})
}
