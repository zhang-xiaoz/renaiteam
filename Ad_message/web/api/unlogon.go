package api

import (
	"context"
	"fmt"
	"image/color"
	"time"
	"web/log"
	"web/models"
	"web/proto_user"

	"web/tool"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mojocn/base64Captcha"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Get_Code_Message(ctx *gin.Context) {
	models.Store = base64Captcha.NewMemoryStore(10240, 1*time.Minute) //实体化一个对象//默认十分钟过期//实体化一个对象
	var dirver base64Captcha.Driver                                   //其实就是一个驱动,到时候驱动DriverString生成图片
	dirvering := base64Captcha.DriverString{                          //封装着属性
		Height:          40,                           //高度
		Width:           120,                          //宽度
		NoiseCount:      0,                            //干扰数
		ShowLineOptions: 4 | 3,                        //展示线条数量
		Length:          5,                            //长度
		Source:          "qwertyuiopasdfghjklzxcvbnm", //验证码随机字符串来源
		BgColor: &color.RGBA{ // 背景颜色
			R: 128,
			G: 128,
			B: 128,
			A: 128,
		},
		Fonts: []string{"wqy-microhei.ttc"}, // 字体
	}
	//按名称加载字符
	dirver = dirvering.ConvertFonts()
	//生成验证码
	c := base64Captcha.NewCaptcha(dirver, models.Store)
	id, b64s, _, err := c.Generate() //生成id  图像和错误
	if err != nil {
		log.SugarLogger.Panic("错误原因:", err) //错误处理简陋
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "失败",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":    "200",
		"id":      id,
		"b64s":    b64s,
		"message": "成功",
	})
}

func Check_Logon(ctx *gin.Context) { //开始验证登录的内容
	//此时验证码已经验证成功
	//获取数据
	login := models.User_login{}                        //实体化一个对象
	err := ctx.ShouldBindBodyWith(&login, binding.JSON) //这是一个通用方法
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, gin.H{
			"code":    "200",
			"message": "验证码不正确",
		})
		return
	}
	//开始验证密码
	if login.Mailbox != "123456" {
		//账户不正确
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "账户不正确",
		})
		return
	}
	if login.Password != "123456" { //密码不正确
		//账户不正确
		ctx.JSON(200, gin.H{
			"code":    "202",
			"message": "密码不正确",
		})
		return
	}
	fmt.Println("aaaaaaa")
	//生成jwt进行登录
	//把token生成好  使用jwt
	jwtjwt := tool.NewJWT()
	//对models签名
	claims1 := models.CustomClaims{
		Mailbox: login.Mailbox,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),         //生效时间
			ExpiresAt: time.Now().Unix() + 60*60, //过期时间 先设置时间短些
		},
	}
	claims2 := models.CustomClaims{
		Mailbox: login.Mailbox,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),            //生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24, //过期时间 先设置时间短些
		},
	}
	access_token, err := jwtjwt.CreateToken(claims1) //短期有效//进行加密
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "后台服务错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	refresh_token, err := jwtjwt.CreateToken(claims2) //短期有效//进行加密
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "后台服务错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_User_Grpc.UserGrpcserver.Host + models.Overall_Situation_User_Grpc.UserGrpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "后台服务错误",
		})
		return
	}
	defer conn.Close()
	userWebClient := proto_user.NewUsersClient(conn)
	//发送给redis进行保存
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto_user.Redis_Storage{
		Key:   login.Mailbox + "access_token",
		Value: access_token,
		Who:   "1",
		Time:  60 * 60,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "后台服务错误",
		})
		return
	}
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto_user.Redis_Storage{
		Key:   login.Mailbox + "refresh_token",
		Value: refresh_token,
		Who:   "1",
		Time:  60 * 60 * 24,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "后台服务错误",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":          "204",
		"access_token":  access_token,
		"refresh_token": refresh_token,
		"message":       "登录成功",
		"uuid":          login.Mailbox,
	})
}
