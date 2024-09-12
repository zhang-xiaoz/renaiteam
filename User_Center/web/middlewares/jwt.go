package middlewares

import (
	"context"
	"encoding/json"
	"time"
	"web/log"
	"web/models"
	modelss "web/models/user"
	"web/proto"
	"web/tool"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 验证jwt两个token//在验证第一个token的时候要不要延期//不要
func JWT_Check() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtt := modelss.JWT_Check{}
		err := ctx.ShouldBindBodyWith(&jwtt, binding.JSON) //这是一个通用方法
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "登录已经过期,请从新登录",
			})
			ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
			return
		}
		//连接Grpc服务
		conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
		//查看连接是否出现错误
		if err != nil {
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "登录已经过期,请从新登录",
			})
			log.SugarLogger.Error(err)
			ctx.Abort()
			return
		}
		defer conn.Close()
		userWebClient := proto.NewUsersClient(conn)
		//开始验证
		if jwtt.Access_token != "" { //这个是短期验证码
			//检验是否存在
			result, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key: jwtt.UUID + "access_token",
				Who: "1",
			})
			if err != nil {
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "短期验证已经过期",
				})
				log.SugarLogger.Error(err)
				ctx.Abort()
				return
			}
			if result.Value != jwtt.Access_token { //验证数据不相等
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				log.SugarLogger.Error(err)
				ctx.Abort()
				return
			}
			//如果短期有效则不用续期
			ctx.Next()
		} else if jwtt.Refresh_token != "" { //这个是长期验证码
			//检验是否存在
			result, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key: jwtt.UUID + "refresh_token",
				Who: "1",
			})
			if err != nil {
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				log.SugarLogger.Error(err)
				ctx.Abort()
				return
			}
			if result.Value != jwtt.Refresh_token { //验证数据不相等
				log.SugarLogger.Error(result.Value)
				log.SugarLogger.Error(jwtt.Refresh_token)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				log.SugarLogger.Error(err)
				ctx.Abort()
				return
			}
			//如果短期无效//长期有效则从新续期//短期从新续期//长期也从新续期
			//查找mailbox是多少
			result1, err := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
				Mailbox: jwtt.UUID,
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			//从新生成两个token
			jwtjwt := tool.NewJWT()
			//对models签名
			claims1 := modelss.CustomClaims{
				Mailbox: result1.Mailbox,
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),         //生效时间
					ExpiresAt: time.Now().Unix() + 60*60, //过期时间 先设置时间短些
				},
			}
			claims2 := modelss.CustomClaims{
				Mailbox: result1.Mailbox,
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),              //生效时间
					ExpiresAt: time.Now().Unix() + 60*60*24*2, //过期时间 先设置时间短些
				},
			}
			access_token, err := jwtjwt.CreateToken(claims1) //短期有效//进行加密
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.Next()
				return
			}
			refresh_token, err := jwtjwt.CreateToken(claims2) //短期有效//进行加密
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.Next()
				return
			}
			//保存到redis数据库里边
			_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key:   jwtt.UUID + "access_token",
				Value: access_token,
				Who:   "1",
				Time:  60 * 60 * 3,
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key:   jwtt.UUID + "refresh_token",
				Value: refresh_token,
				Who:   "1",
				Time:  60 * 60 * 24 * 2,
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			//添加到ctx里边//后续进行保存
			jwtt.Access_token = access_token
			jwtt.Refresh_token = refresh_token
			ccc, err := json.Marshal(jwtt)
			if err != nil {
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			ctx.Set(gin.BodyBytesKey, ccc)
			ctx.Next()
		} else { //两个都为空
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "登录已经过期,请从新登录",
			})
			ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
			return
		}
	}
}

// 验证jwt两个token//在验证第一个token的时候要不要延期//不要
func JWT_Check_New() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//首先获取其中一个短期验证码
		access_token := ctx.GetHeader("access_token")
		refresh_token := ctx.GetHeader("refresh_token")
		uuid := ctx.GetHeader("uuid")
		//连接grpc服务
		//连接Grpc服务
		conn, err := grpc.Dial((models.Overall_Situation_Grpc_Server.Grpcserver.Host + models.Overall_Situation_Grpc_Server.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
		if err != nil {
			log.SugarLogger.Error(err)
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "登录已经过期,请从新登录",
			})
			ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
			return
		}
		defer conn.Close()
		userWebClient := proto.NewUsersClient(conn)
		if access_token != "" { //假如短期验证码不正确则再次发送请求发送长期验证码
			result, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key: uuid + "access_token",
				Who: "1",
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "短期验证已经过期",
				})
				ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
				return
			}
			if result.Value != access_token { //数据不相等
				ctx.JSON(200, map[string]string{
					"code":    "201",
					"message": "短期验证已经过期",
				})
				ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
				return
			}
			ctx.Next() //短期验证码通过不要续期
			return
		} else if refresh_token != "" { //长期验证码
			result, err := userWebClient.Get_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key: uuid + "refresh_token",
				Who: "1",
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
				return
			}
			if result.Value != refresh_token { //数据不相等
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort() //这个会阻止向内层HandlerFunc流动，但是只是阻止作用，并不会使程序停止，所以要用return函数
				return
			}
			//进行续期
			//查找邮箱
			result1, err := userWebClient.Get_User_Mailbox(context.Background(), &proto.Mailbox{
				Mailbox: uuid,
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			//从新生成两个token
			jwtjwt := tool.NewJWT()
			//对models签名
			claims1 := modelss.CustomClaims{
				Mailbox: result1.Mailbox,
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),         //生效时间
					ExpiresAt: time.Now().Unix() + 60*60, //过期时间 先设置时间短些
				},
			}
			claims2 := modelss.CustomClaims{
				Mailbox: result1.Mailbox,
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),              //生效时间
					ExpiresAt: time.Now().Unix() + 60*60*24*2, //过期时间 先设置时间短些
				},
			}
			access_token1, err := jwtjwt.CreateToken(claims1) //短期有效//进行加密
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.Next()
				return
			}
			refresh_token1, err := jwtjwt.CreateToken(claims2) //短期有效//进行加密
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.Next()
				return
			}
			//保存到redis数据库里边
			_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key:   uuid + "access_token",
				Value: access_token1,
				Who:   "1",
				Time:  60 * 60 * 3,
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			_, err = userWebClient.Save_Redis_Storage(context.Background(), &proto.Redis_Storage{
				Key:   uuid + "refresh_token",
				Value: refresh_token1,
				Who:   "1",
				Time:  60 * 60 * 24 * 2,
			})
			if err != nil {
				log.SugarLogger.Error(err)
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			//添加到ctx里边//后续进行保存
			jwt := modelss.Jwt_Check{
				Access_token:  access_token1,
				Refresh_token: refresh_token1,
			}
			ccc, err := json.Marshal(jwt)
			if err != nil {
				ctx.JSON(200, map[string]string{
					"code":    "200",
					"message": "登录已经过期,请从新登录",
				})
				ctx.Abort()
				return
			}
			ctx.Set("jwt", ccc)
			ctx.Next()
			return
		} else {
			ctx.JSON(200, map[string]string{
				"code":    "200",
				"message": "登录已经过期,请从新登录",
			})
			ctx.Abort()
			return
		}
	}
}
