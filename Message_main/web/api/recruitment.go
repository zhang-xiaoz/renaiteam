package api

import (
	"context"
	"web/log"
	"web/models"
	"web/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Get_Training_Personnel(ctx *gin.Context) {
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
	result1, err := pagemessageClient.GetTrainingPersonnel(context.Background(), &proto.PageMemberPaging{
		P:  0,
		Pn: 0,
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

func Get_Training_Message(ctx *gin.Context) {
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
	result1, err := pagemessageClient.GetTrainingMessage(context.Background(), &emptypb.Empty{})
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

func Get_Training_Time(ctx *gin.Context) {
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
	result1, err := pagemessageClient.GetTrainingTime(context.Background(), &emptypb.Empty{})
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
