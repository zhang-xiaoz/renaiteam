package models

//端口文件配置全局变量
var Overall_Situation_Server *Server_Read = &Server_Read{}

//日志文件配置全局变量
var Overall_Situation_Logger *Logger_Read = &Logger_Read{}

//mysql数据库配置全局变量
var Overall_Situation_Mysql *Mysql_Read = &Mysql_Read{}

//redis数据库配置全局变量
var Overall_Situation_Redis *Redis_Read = &Redis_Read{}
