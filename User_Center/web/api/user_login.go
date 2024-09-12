package api

import (
	"context"
	"fmt"
	"os"
	"time"
	"web/log"
	"web/models"
	modelss "web/models/user"
	"web/proto"
	"web/tool"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func User_Login_Get_User_Message(ctx *gin.Context) { //获取用户基本信息
	//已经通过验证能正常登录
	//先再次获取数据
	jwtt := modelss.JWT_Check{}
	err := ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录已经过期,请从新登录",
		})
		return
	}
	//连接grpc服务获取数据
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录已经过期,请从新登录",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//获取mailbox
	mailbox, err := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
		Mailbox: jwtt.UUID,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录已经过期,请从新登录",
		})
		log.SugarLogger.Error(err)
		return
	}
	//获取需要信息
	message, err := userWebClient.Get_User_Mesaage_Mysql(context.Background(), &proto.Mailbox{
		Mailbox: mailbox.Mailbox,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "登录已经过期,请从新登录",
		})
		log.SugarLogger.Error(err)
		return
	}
	//返回信息
	ctx.JSON(200, map[string]string{
		"code":          "202",
		"message":       "成功获取",
		"uuid":          message.Uuid,
		"mailbox":       mailbox.Mailbox,
		"username":      message.Username,
		"name":          message.Name,
		"sex":           message.Sex,
		"address":       message.Address,
		"grade":         message.Grade,
		"direction":     message.Direction,
		"qq":            message.Qq,
		"wechat":        message.Wechat,
		"position":      message.Position,
		"motto":         message.Motto,
		"creat_time":    message.CreatTime,
		"delete_time":   message.DeleteTime,
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Change_User_Message(ctx *gin.Context) { //修改用户基本信息
	//获取数据
	change1 := modelss.JWT_Check{}
	change1.Change_User = make(map[string]string)
	err := ctx.ShouldBindBodyWith(&change1, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败,请从新尝试",
		})
		return
	}
	change2 := proto.User_Change{}
	change2.UserMessage = make(map[string]string)
	//获取能够修改的数据
	for k, v := range change1.Change_User {
		if k == "username" || k == "name" || k == "sex" || k == "direction" || k == "address" || k == "grade" || k == "qq" || k == "wechat" || k == "motto" {
			change2.UserMessage[k] = v
		}
	}
	//连接grpc服务进行更改
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败,请从新尝试",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//要是想成功修改得添加一个mailbox
	//根据uuid获取mailbox
	result6, err := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
		Mailbox: change1.UUID,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败,请从新尝试",
		})
		log.SugarLogger.Error(err)
		return
	}
	change2.UserMessage["mailbox"] = result6.Mailbox
	//更改
	_, err = userWebClient.Revise_User_Message_Mysql(context.Background(), &change2)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "修改失败,请从新尝试",
		})
		log.SugarLogger.Error(err)
		return
	}
	ctx.JSON(200, map[string]string{
		"code":          "203",
		"message":       "修改成功",
		"access_token":  change1.Access_token,
		"refresh_token": change1.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Change_User_Mailbox_Code(ctx *gin.Context) { //修改邮箱
	//获取数据
	jwtt := modelss.JWT_Check{}
	err := ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "验证码错误",
		})
		return
	}
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "获取失败",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//查看验证码是否正确
	//获取验证码数据
	result1, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Key: jwtt.UUID + "_change_mailbox",
		Who: "2",
	})
	//验证得到的验证码是否验证成功
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "验证码未获取",
		})
		return
	}
	if result1.Value != jwtt.Code {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "验证码错误",
		})
		return
	}
	//删除老的验证码
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Key: jwtt.UUID + "_change_mailbox",
		Who: "2",
	})
	//获取新的验证码
	code := tool.Get_Rand_Code(8)
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who:   "1",
		Key:   jwtt.UUID + "修改邮箱验证码通过",
		Time:  60 * 60,
		Value: code,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "验证码未获取",
		})
		return
	}
	//返回数据
	ctx.JSON(200, map[string]string{
		"code":          "205",
		"message":       "验证成功",
		"captcha":       code,
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Get_User_Malibox_Code(ctx *gin.Context) {
	//这里直接获取验证码进行，其他的交给后边的那个判断比较好
	//获取验证码
	jwtt := modelss.JWT_Check{}
	ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "获取失败",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//查看验证码是否已经获取
	backkk, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: jwtt.Mailbox + "登录后修改邮箱",
	})
	if err != nil {
		e, ok := status.FromError(err) //判断有错误是否因为未拥有数据
		if ok {
			switch e.Code() {
			case codes.NotFound:
			default:
				ctx.JSON(200, map[string]string{
					"code":    "203",
					"message": "获取失败",
				})
				return
			}
		}
	}
	//已经有数据
	if backkk != nil {
		if backkk.Time >= 60*4 { //就说明还没有过60秒
			ctx.JSON(200, map[string]string{
				"code":    "205",
				"message": "验证码已经发送",
			})
			return
		}
	}
	code := tool.Get_Rand_Code(8)
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who:   "2",
		Key:   jwtt.Mailbox + "登录后修改邮箱",
		Time:  60 * 5,
		Value: code,
	})
	if err != nil { //出现错误
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "获取失败",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":          "204",
		"message":       "获取成功",
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Change_User_Mailbox_Revise(ctx *gin.Context) { //真正修改邮箱
	//获取数据//新邮箱//验证码code//验证码reason
	//首先根据邮箱判断是否能够使用
	//其次判断验证码
	//最后进行更改
	jwtt := modelss.JWT_Check{}
	ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//查看新邮箱是否有重复
	result, err := userWebClient.Register_Mailbox_Back(context.Background(), &proto.Mailbox{
		Mailbox: jwtt.Mailbox,
	})
	sign := 0 //作为一个标志
	if result != nil {
		if result.Status == 1 {
			ctx.JSON(200, gin.H{
				"code":    "202",
				"message": "该邮箱已经注册",
			})
			return
		} else if result.Status == 2 {
			sign = 1
		} else if result.Status == 3 {
			ctx.JSON(200, gin.H{
				"code":    "202",
				"message": "该用户已经注册",
			})
			return
		} else if result.Status == 4 {
			ctx.JSON(200, gin.H{
				"code":    "202",
				"message": "该用户已经被拉黑",
			})
			return
		}
	}
	if err != nil {
		e, ok := status.FromError(err) //判断有错误是否因为未拥有数据
		if ok {
			switch e.Code() {
			case codes.NotFound:
			default:
				ctx.JSON(200, map[string]string{
					"code":    "203",
					"message": "验证码不正确",
				})
				return
			}
		}
	}
	//进行两个验证码判断
	//判断reason
	//获取验证码数据
	result1, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "1",
		Key: jwtt.UUID + "修改邮箱验证码通过",
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	if result1.Value != jwtt.Reason {
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "修改失败",
		})
		return
	}
	//判断code
	result2, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: jwtt.Mailbox + "登录后修改邮箱",
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	if result2.Value != jwtt.Code {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	//加锁开始进行修改
	//获取原来的邮箱
	result3, err := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
		Mailbox: jwtt.UUID,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "修改失败",
		})
		return
	}
	//进行删除mysql里边的数据
	if sign == 1 {
		_, err = userWebClient.Delete_Mysql_Mailbox(context.Background(), &proto.Mailbox{
			Mailbox: jwtt.Mailbox,
		})
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "204",
				"message": "修改失败",
			})
			return
		}
	}
	//加锁更改用户密码(用和登录一样的锁)(加锁是为了登录时有人修改密码问题)
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                              //用redsync把client封装起来
	rs := redsync.New(pool)                                                                      //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Login + result3.Mailbox) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                         //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	quit := make(chan bool)
	go func() { //go线程
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(time.Second * 4)
				mutex.Extend() //启动一个协程  重置到锁之前状态 保证一直运行
			}
		}
	}()
	defer mutex.Unlock() //解锁//defer是反着执行的
	defer close(quit)
	_, err = userWebClient.Revise_User_Mailbox_Mysql(context.Background(), &proto.Mailbox_Back{
		Mailbox: jwtt.UUID,
		Uuid:    jwtt.Mailbox,
	})
	if err != nil { //有数据
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "修改失败",
		})
		return
	}
	//删除redis里边的数据
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "1",
		Key: jwtt.UUID + "修改邮箱验证码通过",
	})
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: jwtt.Mailbox + "登录后修改邮箱",
	})
	//修改成功
	ctx.JSON(200, map[string]string{
		"code":          "205",
		"message":       "修改成功",
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Change_User_Password(ctx *gin.Context) { //修改密码//一个老密码//一个新密码
	jwtt := modelss.JWT_Check{}
	err := ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		return
	}
	//连接grpc服务
	//连接grpc服务进行更改
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//先验证老密码是否正确
	result4, err := userWebClient.Get_User_Password_Mysql(context.Background(), &proto.Mailbox{
		Mailbox: jwtt.UUID,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	//先验证密码正确不
	result5 := tool.CHECK_Password(jwtt.OldPassword, result4.Mailbox)
	if !result5 {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		return
	}
	newpassword := tool.Salt_Encryption(jwtt.NewPassword) //新密码加密
	//比较新密码和老密码是否一样
	if jwtt.NewPassword == jwtt.OldPassword {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "两个密码一样",
		})
		return
	}
	//获取原来的邮箱
	result10, _ := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
		Mailbox: jwtt.UUID,
	})
	//修改密码//直接进行修改//加锁//加登录的锁
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                               //用redsync把client封装起来
	rs := redsync.New(pool)                                                                       //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Login + result10.Mailbox) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                          //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		return
	}
	quit := make(chan bool)
	go func() { //go线程
		for {
			select {
			case <-quit:
				return
			default:
				time.Sleep(time.Second * 4)
				mutex.Extend() //启动一个协程  重置到锁之前状态 保证一直运行
				log.SugarLogger.Info("注册时redis锁重置8秒")
			}
		}
	}()
	defer mutex.Unlock() //解锁//defer是反着执行的
	defer close(quit)
	//修改密码//修改前加密
	_, err = userWebClient.Revise_User_Password_Mysql(context.Background(), &proto.User_Password{
		Uuid:     jwtt.UUID,
		Password: newpassword,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	ctx.JSON(200, map[string]string{
		"code":          "204",
		"message":       "密码修改成功",
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Change_User_Code(ctx *gin.Context) { //登陆后修改数据获取验证码
	//获取数据
	jwtt := modelss.JWT_Check{}
	ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
	//这个用的是获取url参数
	//格式一般是/Change_User_Code?who=mailbox
	who := ctx.DefaultQuery("who", "无")
	if who != "mailbox" && who != "logout" {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "不能获取url参数错误",
		})
		return
	}
	//发送验证码
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "获取失败",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//查看验证码是否已经获取
	backkk, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: jwtt.UUID + "_change_" + who,
	})
	if err != nil {
		e, ok := status.FromError(err) //判断有错误是否因为未拥有数据
		if ok {
			switch e.Code() {
			case codes.NotFound:
			default:
				ctx.JSON(200, map[string]string{
					"code":    "203",
					"message": "获取失败",
				})
				return
			}
		}
	}
	//已经有数据
	if backkk != nil {
		if backkk.Time >= 60*4 { //就说明还没有过60秒
			ctx.JSON(200, map[string]string{
				"code":    "205",
				"message": "验证码已经发送",
			})
			return
		}
	}
	//开始判断邮箱和uuid是否一致
	result3, err := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
		Mailbox: jwtt.UUID,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "获取失败",
		})
		return
	}
	if result3.Mailbox != jwtt.Mailbox {
		ctx.JSON(200, map[string]string{
			"code":    "206",
			"message": "用户输入不正确",
		})
		return
	}
	code := tool.Get_Rand_Code(8)
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who:   "2",
		Key:   jwtt.UUID + "_change_" + who,
		Time:  60 * 5,
		Value: code,
	})
	if err != nil { //出现错误
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "获取失败",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":          "204",
		"message":       "获取成功",
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作

	})
}

func User_Login_Logout_User(ctx *gin.Context) { //注销账户//只要验证成功就可以进行注销
	//实体化对象
	jwtt := modelss.JWT_Check{}
	err := ctx.ShouldBindBodyWith(&jwtt, binding.JSON)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "注销失败,请从新尝试",
		})
		return
	}
	//连接grpc服务
	//连接grpc服务进行更改
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "注销失败,请从新尝试",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	//获取验证码数据
	result1, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Key: jwtt.UUID + "_change_logout",
		Who: "2",
	})
	//验证得到的验证码是否验证成功
	if err != nil {
		log.SugarLogger.Error(jwtt.UUID + "_change_logout")
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "注销失败,请从新尝试",
		})
		return
	}
	if result1.Value != jwtt.Code {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "验证码错误",
		})
		return
	}
	if len(jwtt.Reason) > 1000 {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "理由太大",
		})
		return
	}
	//删除无关信息//进行退出操作
	//先将状态修改
	_, err = userWebClient.Revise_User_Status_Mysql(context.Background(), &proto.Mailbox_Back{
		Uuid:   jwtt.UUID,
		Status: 2,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "注销失败,请从新尝试",
		})
		return
	}
	//注销理由存储到mysql的password里边
	userWebClient.Revise_User_Password_Mysql(context.Background(), &proto.User_Password{
		Uuid:     jwtt.UUID,
		Password: jwtt.Reason,
	})
	//删除redis验证码
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Key: jwtt.UUID + "_change_logout",
		Who: "1",
	})
	//注销成功
	ctx.JSON(200, map[string]string{
		"code":          "205",
		"message":       "注销成功",
		"access_token":  jwtt.Access_token,
		"refresh_token": jwtt.Refresh_token, //两个都不为空时才进行更改操作
	})
}

func User_Login_Get_User_Picture(ctx *gin.Context) { //获取头像
	// 获取文件名
	filename := ctx.Param("filename")
	fmt.Println(filename)
	// 图片所在的文件夹路径
	dir := "./img/user_picture/"
	// 完整的图片路径
	filePath := dir + filename
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "获取失败",
		})
		return
	}
	// 使用Gin的c.File方法发送文件
	ctx.File(filePath)
}

func User_Login_Change_User_Picture(ctx *gin.Context) { //上传图片
	//同意某个图片
	//首先判断get里边是否有数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	// 获取文件名
	filename := ctx.GetHeader("uuid")
	if filename == "" {
		ctx.JSON(200, map[string]string{
			"code":    "206",
			"message": "用户不是有效身份",
			"jwt":     jwtString,
		})
		return
	}
	file, error := ctx.FormFile("file")
	if error != nil {
		log.SugarLogger.Error(error)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	if file == nil {
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	if file.Size > 2*1024*1024 {
		//返回信息
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "图片过大",
			"jwt":     jwtString,
		})
		return
	}
	contentType := file.Header.Get("Content-Type")
	attribute := ""
	if contentType == "image/jpeg" {
		attribute = ".jpg"
	} else if contentType == "image/png" {
		attribute = ".png"
	} else {
		//返回信息
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "图片格式不对",
			"jwt":     jwtString,
		})
		return
	}
	//先删除图片
	os.Remove("./img/user_picture/" + filename + ".jpg")
	os.Remove("./img/user_picture/" + filename + ".png")
	//保存图片
	err := ctx.SaveUploadedFile(file, "./img/user_picture/"+filename+attribute)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "上传失败",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "205",
		"message": "上传成功",
		"jwt":     jwtString,
	})
}

func Get_User_Picture_Name(ctx *gin.Context) { //获取图片名字
	//同意某个图片
	//首先判断get里边是否有数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	// 获取文件名
	filename := ctx.GetHeader("uuid")
	if filename == "" {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "用户不是有效身份",
			"jwt":     jwtString,
		})
		return
	}
	//开始获取图片名字
	dir := "./img/user_picture/"
	// 完整的图片路径
	filePath := dir + filename + ".jpg"
	if _, err1 := os.Stat(filePath); os.IsNotExist(err1) {
		//错误
		filePath = filename + ".png"
		if _, err2 := os.Stat(filePath); os.IsNotExist(err2) {
			ctx.JSON(200, map[string]string{
				"code":    "203",
				"message": "默认.png",
				"jwt":     jwtString,
			})
			return

		}
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": filename + ".png",
			"jwt":     jwtString,
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": filename + ".jpg",
		"jwt":     jwtString,
	})
}

func User_Quit(ctx *gin.Context) { //退出登录
	jwtt := modelss.JWT_Check{}
	err := ctx.ShouldBind(&jwtt)
	if err != nil {
		log.SugarLogger.Error(err)
		return
	}
	//开始退出登录
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto.NewUsersClient(conn)
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "1",
		Key: jwtt.UUID + "refresh_token",
	})
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "1",
		Key: jwtt.UUID + "access_token",
	})
}
