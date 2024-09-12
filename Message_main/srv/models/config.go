package models

//端口文件配置全局变量
var Overall_Situation_Server *Server_Read = &Server_Read{}

//日志文件配置全局变量
var Overall_Situation_Logger *Logger_Read = &Logger_Read{}

//mysql数据库配置全局变量
var Overall_Situation_Mysql *Mysql_Read = &Mysql_Read{}

//redis数据库配置全局变量
var Overall_Situation_Redis *Redis_Read = &Redis_Read{}

//读取配置文件
type Server_Read struct {
	Server Server `yaml:"server"`
}

//端口配置文件
type Server struct {
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

//读取mysql配置文件
type Mysql_Read struct {
	Mysql Mysql `yaml:"mysql"`
}

//mysql配置文件内容
type Mysql struct {
	User     string `yaml:"user"`     //用户名
	Password string `yaml:"password"` //用户密码
	Port     string `yaml:"port"`     //端口
	Name     string `yaml:"name"`     //数据库名
}

//读取redis配置文件
type Redis_Read struct {
	Redis Redis `yaml:"redis"`
}

//redis配置文件内容
type Redis struct {
	Port     string `yaml:"port"`     //端口
	Password string `yaml:"password"` //密码
}
