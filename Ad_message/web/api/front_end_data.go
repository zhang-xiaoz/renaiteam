package api

import (
	"context"
	"web/log"
	"web/models"
	"web/proto_page"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Revise_Message_Mysql(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Message_Mysql{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	//开始修改数据关于我们的value1值
	_, err = userWebClient.Revise_Message(context.Background(), &proto_page.ReviseMessage{
		Key:     message1.Key,
		Keyform: message1.Keyform,
		Value1:  message1.Value1,
		Value2:  message1.Value2,
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
		"message": "修改成功",
		"jwt":     jwtString,
	})
}

func Revise_Training_Time(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Message_Mysql{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	//开始修改数据关于我们的value1值
	_, err = userWebClient.Revise_Message(context.Background(), &proto_page.ReviseMessage{
		Key:     message1.Key,
		Keyform: "培训时间",
		Value1:  message1.Value1,
		Value2:  message1.Value2,
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
		"message": "修改成功",
		"jwt":     jwtString,
	})
}

func Add_Training_Time(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Add_Training_Time{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	//开始修改数据关于我们的value1值
	_, err = userWebClient.Add_Training_Time(context.Background(), &proto_page.ReviseMessage{
		Key:     message1.Key,
		Keyform: "",
		Value1:  message1.Value1,
		Value2:  message1.Value2,
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

func Delete_Training_Time(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取所有数据
	message1 := []models.Message_Mysql{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	//开始转换数据
	resultt := []*proto_page.ReviseMessage{}
	for _, m := range message1 {
		resulttt := proto_page.ReviseMessage{
			Key:     m.Key,
			Keyform: m.Keyform,
			Value1:  m.Value1,
			Value2:  m.Value2,
		}
		resultt = append(resultt, &resulttt)
	}
	_, err = userWebClient.Del_Training_TIme(context.Background(), &proto_page.DelMessage{
		Message: resultt,
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

func Add_Club_Direction(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取需要修改的数据
	message1 := models.Message_Mysql{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	//开始修改数据关于我们的value1值
	_, err = userWebClient.Add_Club_Direction(context.Background(), &proto_page.ReviseMessage{
		Key:     message1.Key,
		Keyform: message1.Keyform,
		Value1:  message1.Value1,
		Value2:  message1.Value2,
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

func Delete_Club_Direction(ctx *gin.Context) {
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取删除数据
	message1 := models.Message_Mysql{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	_, err = userWebClient.Del_Club_Direction(context.Background(), &proto_page.ReviseMessage{
		Key: message1.Key,
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

func Revise_Award_Information(ctx *gin.Context) {
	//修改获奖信息//根据uuid
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取删除数据
	message1 := models.Prize{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	_, err = userWebClient.Revise_Award_Information(context.Background(), &proto_page.Prize{
		Uuid:   message1.Uuid,
		Name:   message1.Name,
		Adards: message1.Awards,
		Time:   message1.Time,
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
		"message": "修改成功",
		"jwt":     jwtString,
	})
}

func Delete_Award_Information(ctx *gin.Context) {
	//修改获奖信息//根据uuid
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取删除数据
	message1 := models.Prize{}
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
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	_, err = userWebClient.Del_Award_Information(context.Background(), &proto_page.Prize{
		Uuid: message1.Uuid,
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

func Add_Award_Information(ctx *gin.Context) {
	//修改获奖信息//根据uuid
	//开始获取jwt数据
	jwtString := ""
	jwt, _ := ctx.Get("jwt")
	if jwt != nil {
		jwtString, _ = jwt.(string)
	}
	//开始获取删除数据
	message1 := models.Prize{}
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
	message1.Uuid = uuid.New().String()
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
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
	userWebClient := proto_page.NewPagemessageClient(conn)
	_, err = userWebClient.Add_Award_Information(context.Background(), &proto_page.Prize{
		Uuid:   message1.Uuid,
		Name:   message1.Name,
		Adards: message1.Awards,
		Time:   message1.Time,
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
