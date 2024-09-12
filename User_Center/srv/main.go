package main

import (
	"net"
	"srv/config"
	"srv/hander"
	"srv/initialization"
	"srv/log"
	"srv/proto"

	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()        //初始化配置文件
	log.InitLogger()           //初始化日志文件
	initialization.InitMysql() //初始化mysql
	initialization.InitRedis() //初始化redis
	server := grpc.NewServer()
	proto.RegisterUsersServer(server, &hander.UserServer{})
	lis, _ := net.Listen("tcp", ":50051")
	log.SugarLogger.Info("服务器端口开启")
	server.Serve(lis)
}
