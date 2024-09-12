package models

import "github.com/mojocn/base64Captcha"

//读取配置文件结构体

//日志文件配置全局变量
var Overall_Situation_Logger *Logger_Read = &Logger_Read{}

// 端口文件配置全局变量
var Overall_Situation_Server *Server_Read = &Server_Read{}

//grpc服务
var Overall_Situation_Grpc *Grpc_Server_Read = &Grpc_Server_Read{}

//grpc服务读取配置文件
type Grpc_Server_Read struct {
	Grpcserver Grpcserver `yaml:"grpcserver"`
}

//usergrpc服务端口配置文件
type Grpcserver struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

//grpc_user服务
var Overall_Situation_User_Grpc *User_Grpc_Server_Read = &User_Grpc_Server_Read{}

//grpc服务读取配置文件
type User_Grpc_Server_Read struct {
	UserGrpcserver UserGrpcserver `yaml:"usergrpcserver"`
}

//usergrpc服务端口配置文件
type BlogGrpcserver struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

//grpc_user服务
var Overall_Situation_Blog_Grpc *Blog_Grpc_Server_Read = &Blog_Grpc_Server_Read{}

//grpc服务读取配置文件
type Blog_Grpc_Server_Read struct {
	BlogGrpcserver UserGrpcserver `yaml:"bloggrpcserver"`
}

//usergrpc服务端口配置文件
type UserGrpcserver struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

//读取日志文件配置文件
type Logger_Read struct {
	Logger Logger `yaml:"logger"`
}

// 日志文件配置文件内容
type Logger struct {
	Filename   string `yaml:"filename"`  // 文件位置
	MaxSize    int    `yaml:"maxsize"`   // 进行切割之前,日志文件的最大大小(MB为单位)
	MaxSave    int    `yaml:"maxsave"`   // 保留旧文件的最大天数
	MaxNumbers int    `yaml:"maxNumber"` // 保留旧文件的最大个数
	Compress   bool   `yaml:"compress"`  // 是否压缩/归档旧文件
}

//读取配置文件
type Server_Read struct {
	Server Server `yaml:"server"`
}

//端口配置文件
type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var Store base64Captcha.Store //图片验证码
