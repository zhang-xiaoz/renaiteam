package api

import (
	"context"
	"image/color"
	"time"
	"web/log"
	"web/models"
	modelss "web/models/user"
	"web/proto"
	"web/tool"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	goredislib "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// 未登录时注册时邮箱验证码获取
func User_Unlisted_Get_Mailbox_Code(ctx *gin.Context) {
	mailbox := modelss.Mailbox_Dode{}
	//获取邮箱
	err := ctx.ShouldBind(&mailbox)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "邮箱格式填写错误",
		})
		return
	}
	//查看是否已经注册//未注册//已经注册//等待管理员验证
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "邮箱格式填写错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	//生成grpc的client并调用接口
	userWebClient := proto.NewUsersClient(conn)
	//获取数据进行判断
	get1, err := userWebClient.Register_Mailbox_Back(context.Background(), &proto.Mailbox{
		Mailbox: mailbox.Mailbox,
	})
	a := 0
	if err != nil {
		//查看是否是因为数据为空造成
		st, ok := status.FromError(err)
		if ok {
			// 检查状态代码是否为 codes.NotFound
			if st.Code() == codes.NotFound {
				a = 1
			}
		}
		//现在判断出不是为空导致错误
		if a == 0 {
			ctx.JSON(200, map[string]string{
				"code":    "201",
				"message": "邮箱格式填写错误",
			})
			return
		}
	}
	//没有错误获取到数据且不为空
	if a == 0 {
		if get1.Status == 1 {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "正常用户已经完成过注册",
			})
			return
		}
		if get1.Status == 3 {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "已经完成注册请等待审核",
			})
			return
		}
		if get1.Status == 4 {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "已经被拉黑不允许注册",
			})
			return
		}
	}
	//用户未注册或者注销过一次的用户//开始进行下边判断
	//验证码是否已经获取
	backkk, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: mailbox.Mailbox + "注册",
	})
	if err != nil {
		e, ok := status.FromError(err) //判断有错误是否因为未拥有数据
		if ok {
			switch e.Code() {
			case codes.NotFound:
			default:
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "邮箱格式填写错误",
				})
				return
			}
		}
	}
	//已经有数据
	if backkk != nil {
		if backkk.Time >= 60*4 { //就说明还没有过60秒
			ctx.JSON(200, map[string]string{
				"code":    "203",
				"message": "验证码已经发送",
			})
			return
		}
	}
	//没有数据或者可以再次获取数据//验证码发送
	code := tool.Get_Rand_Code(8)
	//进行存储到redis2里边
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who:   "2",
		Key:   mailbox.Mailbox + "注册",
		Time:  60 * 5,
		Value: code,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "邮箱格式填写错误",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "204",
		"message": "发送成功",
	})
}

// 未登录时进行注册//进行验证
func User_Unlisted_Register_Mailbox(ctx *gin.Context) {
	register := modelss.Mailbox_Register{}
	//获取信息
	err := ctx.ShouldBind(&register)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "格式填写错误",
		})
		return
	}
	//将密码进行盐值加密
	//生成唯一表示uuid//这里我选用邮箱md5加密生成的东西作为uuid
	newpassword := tool.Salt_Encryption(register.Password)
	//进行加锁，加锁是为了避免查找是显示不存在，但是两个用户同时去注册会一样成功，所以最后添加信息时进行加锁
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                                  //用redsync把client封装起来
	rs := redsync.New(pool)                                                                          //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Register + register.Mailbox) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                             //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "格式填写错误",
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
	//查看是否已经注册//未注册//已经注册//等待管理员验证
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "格式填写错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	//生成grpc的client并调用接口
	userWebClient := proto.NewUsersClient(conn)
	//获取数据进行判断
	get1, err := userWebClient.Register_Mailbox_Back(context.Background(), &proto.Mailbox{
		Mailbox: register.Mailbox,
	})
	a := 0
	if err != nil {
		//查看是否是因为数据为空造成
		st, ok := status.FromError(err)
		if ok {
			// 检查状态代码是否为 codes.NotFound
			if st.Code() == codes.NotFound {
				a = 1
			}
		}
		//现在判断出不是为空导致错误
		if a == 0 {
			ctx.JSON(200, map[string]string{
				"code":    "201",
				"message": "格式填写错误",
			})
			return
		}
	}
	//没有错误获取到数据且不为空
	if a == 0 {
		if get1.Status == 1 {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "正常用户已经完成过注册",
			})
			return
		}
		if get1.Status == 3 {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "已经完成注册请等待审核",
			})
			return
		}
		if get1.Status == 4 {
			ctx.JSON(200, map[string]string{
				"code":    "202",
				"message": "已经被拉黑不允许注册",
			})
			return
		}
	}
	//验证是否获取验证码
	backkk, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: register.Mailbox + "注册",
	})
	if err != nil {
		//如果出现错误//首先可能是真出现问题或者是没有获取到数据
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "未获取邮箱验证码",
		})
		return
	}
	if backkk == nil {
		ctx.JSON(200, map[string]string{
			"code":    "203",
			"message": "未获取邮箱验证码",
		})
		return
	}
	//验证验证码是否正确
	if backkk.Value != register.Code { //不正确
		ctx.JSON(200, map[string]string{
			"code":    "204",
			"message": "验证码不正确",
		})
		return
	}
	//删除redis里边数据
	userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: register.Mailbox + "注册",
	})
	//进行存储//直接存储就行
	userWebClient.Add_User_Message(context.Background(), &proto.User_Message{
		Status:   3,
		Uuid:     uuid.New().String(),
		Mailbox:  register.Mailbox,
		Password: newpassword,
		Username: register.Name,
		Name:     register.Name,
		Sex:      "未知",
		Address:  "未知",
	})
	//添加判断
	ctx.JSON(200, map[string]string{
		"code":    "205",
		"message": "注册成功,请等待审核",
	})
}

// 未登录时获取图片验证
func User_Unlisted_Get_Picture_Code(ctx *gin.Context) {
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

// 未登录时进行登录
func User_Unlisted_Login(ctx *gin.Context) {
	//此时验证码已经验证成功
	//获取数据
	login := modelss.User_login{}                       //实体化一个对象
	err := ctx.ShouldBindBodyWith(&login, binding.JSON) //这是一个通用方法
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "204",
			"message": "格式不正确",
		})
		return
	}
	//查找用户是否能够登录
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	//查看连接是否出现错误
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "验证码不正确",
		})
		log.SugarLogger.Error(err)
		return
	}
	defer conn.Close()
	//加锁更改用户密码(用和登录一样的锁)(加锁是为了登录时有人修改密码问题)
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                            //用redsync把client封装起来
	rs := redsync.New(pool)                                                                    //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Login + login.Mailbox) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                       //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
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
				log.SugarLogger.Info("注册时redis锁重置8秒")
			}
		}
	}()
	defer mutex.Unlock() //解锁//defer是反着执行的
	defer close(quit)
	//生成grpc的client并调用接口
	userWebClient := proto.NewUsersClient(conn)
	result, err := userWebClient.Get_User_Mesaage_Mysql(context.Background(), &proto.Mailbox{
		Mailbox: login.Mailbox,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "用户未注册",
		})
		log.SugarLogger.Error(err)
		return
	}
	//1表示正常用户 2表示注销过的用户 3表示未通过审核  4表示改用户不允许注册已经被拉黑
	if result.Status == 2 {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "该用户已经注销",
		})
		return
	} else if result.Status == 3 {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "该用户未通过审核",
		})
		return

	} else if result.Status == 4 {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "该用户已经拉黑不允许注册",
		})
		return
	}
	//已经确定是正常用户
	//验证密码是否正确
	r := tool.CHECK_Password(login.Password, result.Password)
	if !r { //密码错误
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		return
	}
	//生成jwt进行登录
	//把token生成好  使用jwt
	jwtjwt := tool.NewJWT()
	//对models签名
	claims1 := modelss.CustomClaims{
		Mailbox: login.Mailbox,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),         //生效时间
			ExpiresAt: time.Now().Unix() + 60*60, //过期时间 先设置时间短些
		},
	}
	claims2 := modelss.CustomClaims{
		Mailbox: login.Mailbox,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),              //生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*2, //过期时间 先设置时间短些
		},
	}
	access_token, err := jwtjwt.CreateToken(claims1) //短期有效//进行加密
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	refresh_token, err := jwtjwt.CreateToken(claims2) //短期有效//进行加密
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		log.SugarLogger.Error(err)
		return
	}
	//发送给redis进行保存
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Key:   result.Uuid + "access_token",
		Value: access_token,
		Who:   "1",
		Time:  60 * 60 * 3,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		return
	}
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Key:   result.Uuid + "refresh_token",
		Value: refresh_token,
		Who:   "1",
		Time:  60 * 60 * 24 * 2,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "密码错误",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":          "203",
		"access_token":  access_token,
		"refresh_token": refresh_token,
		"message":       "登录成功",
		"uuid":          result.Uuid,
	})
}

// 未登录时验证码获取(用于更改密码)
func User_Unlisted_Get_Code(ctx *gin.Context) {
	//实体化对象
	mailbox := modelss.Mailbox_Dode{}
	err := ctx.ShouldBind(&mailbox)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "用户未注册",
		})
		return
	}
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "用户未注册",
		})
		return
	}
	defer conn.Close()
	//生成client接口
	userWebClient := proto.NewUsersClient(conn)
	result, err := userWebClient.Get_User_Mesaage_Mysql(context.Background(), &proto.Mailbox{
		Mailbox: mailbox.Mailbox,
	})
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "用户未注册",
		})
		log.SugarLogger.Error(err)
		return
	}
	//1表示正常用户 2表示注销过的用户 3表示未通过审核  4表示改用户不允许注册已经被拉黑
	if result.Status == 2 {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "该用户已经注销",
		})
		return
	} else if result.Status == 3 {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "该用户未通过审核",
		})
		return

	} else if result.Status == 4 {
		ctx.JSON(200, map[string]string{
			"code":    "202",
			"message": "该用户已经拉黑不允许注册",
		})
		return
	}
	//已经确定是正常用户
	//先查找验证码是否已经获取
	//验证码是否已经获取
	backkk, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: mailbox.Mailbox + "验证码",
	})
	if err != nil {
		e, ok := status.FromError(err) //判断有错误是否因为未拥有数据
		if ok {
			switch e.Code() {
			case codes.NotFound:
			default:
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "用户未注册",
				})
				return
			}
		}
	}
	//已经有数据
	if backkk != nil {
		if backkk.Time >= 60*4 { //就说明还没有过60秒
			ctx.JSON(200, map[string]string{
				"code":    "204",
				"message": "验证码已经发送",
			})
			return
		}
	}
	//发送验证码
	code := tool.Get_Rand_Code(8)
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who:   "2",
		Key:   mailbox.Mailbox + "验证码",
		Time:  60 * 5,
		Value: code,
	})
	if err != nil { //出现错误
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "用户未注册",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "203",
		"message": "发送成功",
	})
}

// 未登录时验证码验证(用于更改密码)
func User_Unlisted_Check_Code(ctx *gin.Context) {
	//实体化对象(信息有验证码,邮箱)
	check := modelss.User_Code_Check{}
	err := ctx.ShouldBind(&check)
	if err != nil {
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "验证码不正确",
		})
	}
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "验证码不正确",
		})
		return
	}
	defer conn.Close()
	//生成client接口
	userWebClient := proto.NewUsersClient(conn)
	//因为如果邮箱不存在是无法发送验证码的，所以不必在判断是否存在该邮箱
	//获取redis数据
	result, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: check.Mailbox + "验证码",
	})
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "202",
			"message": "验证码未获取",
		})
		return
	}
	if result.Value != check.Code {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	//删除redis数据
	_, err = userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: check.Mailbox + "验证码",
	})
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	//如果通过了这个验证则赋予一个数据证明//添加一个redis进行证明
	code := tool.Get_Rand_Code(8)
	_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who:   "2",
		Key:   check.Mailbox + "验证码通过",
		Time:  60 * 60,
		Value: code,
	})
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "203",
			"message": "验证码不正确",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code":         "204",
		"message":      "验证码正确",
		"mailbox_code": code,
		"mailbox":      check.Mailbox,
	})
}

// 未登录时密码更改(验证码已经通过验证)
func User_Unlisted_Change_Password(ctx *gin.Context) {
	//首先获取数据
	change := modelss.User_Code_Twice_Check{}
	err := ctx.ShouldBind(&change)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "超时",
		})
	}
	//还得改密码格式来
	newpassword := tool.Salt_Encryption(change.Password)
	//连接grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "超时",
		})
		return
	}
	defer conn.Close()
	//生成client接口
	userWebClient := proto.NewUsersClient(conn)
	//获取验证码进行判断
	result, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: change.Mailbox + "验证码通过",
	})
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "超时",
		})
		return
	}
	if result.Value != change.Code {
		ctx.JSON(200, gin.H{
			"code":    "201",
			"message": "超时",
		})
		return
	}
	//通过验证
	//加锁更改用户密码(用和登录一样的锁)(加锁是为了登录时有人修改密码问题)
	client := goredislib.NewClient(&goredislib.Options{ //连接redis服务器配置//redis包内容
		Addr:     models.Overall_Situation_Redis.Redis.Port,
		Password: models.Overall_Situation_Redis.Redis.Password,
	})
	pool := goredis.NewPool(client)                                                             //用redsync把client封装起来
	rs := redsync.New(pool)                                                                     //redsync.New方法从给定的Redis连接池创建并返回一个新的Redsync实例。
	mutex := rs.NewMutex(models.Overall_Situation_Redisclock.RedisClock.Login + change.Mailbox) //创建一个redis锁//里边包含锁的一些配置//默认过期时间八秒
	if err := mutex.Lock(); err != nil {                                                        //获取锁
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "超时",
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
	//更改密码//更改不成功显示超时
	_, err = userWebClient.Revise_User_Message_Mysql(context.Background(), &proto.User_Change{
		UserMessage: map[string]string{
			"password": newpassword,
			"mailbox":  change.Mailbox,
		},
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "201",
			"message": "超时",
		})
		return
	}
	//删除redis里边数据数据
	_, err = userWebClient.Delete_Redis_Storage(context.Background(), &proto.Redis_Storage{
		Who: "2",
		Key: change.Mailbox + "验证码通过",
	})
	if err != nil {
		log.SugarLogger.Error(err)
	}
	ctx.JSON(200, map[string]string{
		"code":    "202",
		"message": "更改成功",
	})
}
