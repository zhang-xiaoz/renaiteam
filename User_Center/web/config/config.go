package config

import (
	"web/models"

	"github.com/spf13/viper"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile("config/config.yaml") //地址
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(models.Overall_Situation_Logger); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(models.Overall_Situation_Server); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(models.Overall_Situation_Grpc_Server); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(models.Overall_Situation_Redis); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(models.Overall_Situation_Redisclock); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(models.JWTconfig); err != nil {
		panic(err)
	}
}
