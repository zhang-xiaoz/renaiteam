package api

import (
	"context"
	"strconv"
	"web/log"
	"web/models"
	"web/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Get_Member_Message(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//获取是获取的哪一个成员
	message := models.Message{}
	err := ctx.ShouldBind(&message)
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	defer conn.Close()
	pagemessageClient := proto.NewPagemessageClient(conn)

	result1, err := pagemessageClient.GetMemberMessage(context.Background(), &proto.PageMemberPaging{
		P:       p,
		Pn:      pn,
		Message: message.Message,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "201",
		"message": "成功",
		"data":    result1.Message,
	})
}

func Get_Prize_Message(ctx *gin.Context) {
	//获取两个参数
	p, _ := strconv.ParseInt(ctx.Query("p"), 10, 64)
	pn, _ := strconv.ParseInt(ctx.Query("pn"), 10, 64)
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	defer conn.Close()
	pagemessageClient := proto.NewPagemessageClient(conn)
	result1, err := pagemessageClient.GetPrizeMessage(context.Background(), &proto.PageMemberPaging{
		P:  p,
		Pn: pn,
	})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "201",
		"message": "成功",
		"data":    result1.Message,
	})
}

func Get_Club_Direction_Message(ctx *gin.Context) {
	//连接Grpc服务
	conn, err := grpc.Dial((models.Overall_Situation_Grpc.Grpcserver.Host + models.Overall_Situation_Grpc.Grpcserver.Port), grpc.WithTransportCredentials(insecure.NewCredentials())) //拨号，建立连接//后边参数安全参数
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	defer conn.Close()
	pagemessageClient := proto.NewPagemessageClient(conn)
	result1, err := pagemessageClient.GetClubDirectionMessage(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.SugarLogger.Error(err)
		ctx.JSON(200, map[string]string{
			"code":    "200",
			"message": "发生错误",
		})
		return
	}
	ctx.JSON(200, map[string]string{
		"code":    "201",
		"message": "成功",
		"data":    result1.Message,
	})
}
