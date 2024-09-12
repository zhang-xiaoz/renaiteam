package api

import (
	"context"
	"fmt"
	"strconv"
	"web/log"
	"web/models"
	"web/proto_user"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Get_Register_User(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找未审核的用户
	result1, err := userWebClient.Get_Register_User(context.Background(), &proto_user.MemberPaging{
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
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result1.Message,
	})
}

func Agree_Register_User(ctx *gin.Context) {
	//获取数据
	login := models.Register{}                          //实体化一个对象
	err := ctx.ShouldBindBodyWith(&login, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, gin.H{
			"code":    "202",
			"message": "发生错误",
		})
		return
	}
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	_, err = userWebClient.Revise_User_Status(context.Background(), &proto_user.User_Revise_Status{
		Status:  1,
		Mailbox: login.Mailbox,
		Grade:   login.Grade,
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
		"message": "成功",
		"jwt":     jwtString,
	})
}

func Refuse_Register_User(ctx *gin.Context) {
	//获取数据
	login := models.Mailbox{}                           //实体化一个对象
	err := ctx.ShouldBindBodyWith(&login, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, gin.H{
			"code":    "202",
			"message": "发生错误",
		})
		return
	}
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	_, err = userWebClient.Refuse_User_Status(context.Background(), &proto_user.Mailbox{
		Mailbox: login.Mailbox,
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
		"message": "成功",
		"jwt":     jwtString,
	})
}

func Get_User(ctx *gin.Context) {
	//获取正常用户
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找未审核的用户
	result1, err := userWebClient.Get_User(context.Background(), &proto_user.MemberPaging{
		P:       p,
		Pn:      pn,
		Message: "正常用户",
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
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result1.Message,
	})
}

func Revise_User(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取数据
	change1 := models.Revise_User{}
	change1.Change_User = make(map[string]string)
	err := ctx.ShouldBindBodyWith(&change1, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败",
			"jwt":     jwtString,
		})
		return
	}
	change2 := proto_user.User_Change{}
	change2.UserMessage = make(map[string]string)
	//获取能够修改的数据
	for k, v := range change1.Change_User {
		if k == "sex" || k == "username" || k == "name" || k == "address" || k == "grade" || k == "direction" || k == "qq" || k == "wechat" || k == "position" || k == "motto" || k == "creat_time" {
			change2.UserMessage[k] = v
		}
	}
	fmt.Println(change1)
	//连接grpc服务进行更改
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败",
			"jwt":     jwtString,
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	userWebClient := proto_user.NewUsersClient(conn)
	change2.UserMessage["mailbox"] = change1.Mailbox
	//更改数据
	_, err = userWebClient.Revise_User_Message_Mysql(context.Background(), &change2)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败",
			"jwt":     jwtString,
		})
		log.SugarLogger.Error(err)
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "修改成功",
		"jwt":     jwtString,
	})
}

func Seek_Grade_User(ctx *gin.Context) {
	//获取正常用户
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取参数
	message := models.Grade{}
	err := ctx.ShouldBindBodyWith(&message, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	result1, err := userWebClient.Seek_Grade_User(context.Background(), &proto_user.MemberPaging{
		P:       p,
		Pn:      pn,
		Message: message.Grade,
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
		"code":    "204",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result1.Message,
	})
}

func Seek_Name_User(ctx *gin.Context) {
	//获取正常用户
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取参数
	message := models.Name{}
	err := ctx.ShouldBindBodyWith(&message, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	result1, err := userWebClient.Seek_Name_User(context.Background(), &proto_user.MemberPaging{
		P:       p,
		Pn:      pn,
		Message: message.Name,
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
		"code":    "204",
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result1.Message,
	})
}

func Get_Cancel_User(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	result1, err := userWebClient.Get_Cancel_User(context.Background(), &proto_user.MemberPaging{
		P:       p,
		Pn:      pn,
		Message: "注销用户",
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
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result1.Message,
	})
}

func Delete_Cancel_User(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取参数
	message := models.Mailbox{}
	err := ctx.ShouldBindBodyWith(&message, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	_, err = userWebClient.Delete_Blacklist_User(context.Background(), &proto_user.Mailbox{
		Mailbox: message.Mailbox,
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

func Get_Blacklist_User(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	result1, err := userWebClient.Get_Blacklist_User(context.Background(), &proto_user.MemberPaging{
		P:       p,
		Pn:      pn,
		Message: "黑名单",
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
		"message": "获取成功",
		"jwt":     jwtString,
		"data":    result1.Message,
	})
}

func Delete_Blacklist_User(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取参数
	message := models.Mailbox{}
	err := ctx.ShouldBindBodyWith(&message, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	_, err = userWebClient.Delete_Blacklist_User(context.Background(), &proto_user.Mailbox{
		Mailbox: message.Mailbox,
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

func Add_Blacklist_User(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取参数
	message := models.Blacklist_User{}
	err := ctx.ShouldBindBodyWith(&message, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	message.Uuid = uuid.New().String()
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	_, err = userWebClient.Add_Blacklist_User(context.Background(), &proto_user.Blacklist_User{
		Mailbox:  message.Mailbox,
		Uuid:     message.Uuid,
		Status:   4,
		Password: message.Password,
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
		"message": "增加成功",
		"jwt":     jwtString,
	})
}

func Delete_Revise_User(ctx *gin.Context) {
	//删除用户注册信息
	//删除用户头像信息
	//删除用户所有文章信息
	//删除用户所有信息  。。。

}

func Add_Register_Blacklist(ctx *gin.Context) {
	//加入黑名单就是把sttatus和password改了
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//获取参数
	message := models.Blacklist_User{}
	err := ctx.ShouldBindBodyWith(&message, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "发生错误",
			"jwt":     jwtString,
		})
		return
	}
	//连接UserGrpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_user.NewUsersClient(conn)
	//开始查找数据
	_, err = userWebClient.Add_Register_Blacklist(context.Background(), &proto_user.Blacklist_User{
		Mailbox:  message.Mailbox,
		Status:   4,
		Password: message.Password,
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
		"code":    "204",
		"message": "拉入成功",
		"jwt":     jwtString,
	})

}
