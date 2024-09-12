package main

import (
	"net"
	"srv/config"
	"srv/hander"
	"srv/initialization"
	"srv/log"
	"srv/models"
	"srv/proto"

	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()        //初始化配置文件
	log.InitLogger()           //初始化日志文件
	initialization.InitMysql() //初始化mysql
	initialization.InitRedis() //初始化redis
	server := grpc.NewServer()
	proto.RegisterPagemessageServer(server, &hander.PagemessageServer{})
	lis, _ := net.Listen("tcp", models.Overall_Situation_Server.Server.Port)
	log.SugarLogger.Info("服务器端口开启")
	server.Serve(lis)
}
