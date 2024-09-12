package models

import (
	"web/models/user"

	"github.com/mojocn/base64Captcha"
)

// 日志文件配置全局变量
var Overall_Situation_Logger *Logger_Read = &Logger_Read{}

// 端口文件配置全局变量
var Overall_Situation_Server *Server_Read = &Server_Read{}

var Overall_Situation_Grpc_Server *Grpc_Server_Read = &Grpc_Server_Read{}

// redis数据库配置全局变量
var Overall_Situation_Redis *Redis_Read = &Redis_Read{}

var Overall_Situation_Redisclock *Redis_Clock = &Redis_Clock{}

var Store base64Captcha.Store //图片验证码

var JWTconfig *user.JWTConfig = &user.JWTConfig{} //jwt保密验证信息
