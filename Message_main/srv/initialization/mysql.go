package initialization

import (
	"srv/log"
	"srv/models"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 连数据库gorm方法
var (
	DB  *gorm.DB
	Err error
)

func InitMysql() {
	go func() {
		dsn := models.Overall_Situation_Mysql.Mysql.User + ":" + models.Overall_Situation_Mysql.Mysql.Password + "@tcp(" + models.Overall_Situation_Mysql.Mysql.Port + ")/" + models.Overall_Situation_Mysql.Mysql.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
		DB, Err = gorm.Open(mysql.Open(dsn), &gorm.Config{}) //gorm.Config{}里边可以加入配置
		if Err != nil {
			log.SugarLogger.Error(Err)
		} else {
			sql, _ := DB.DB() //返回sqlDb //为什么在创建连接数时不直接在gorm中创建而要调用mysql中的Db//网上一些博主的观点//创建gorm.DB对象的时候连接并没有被创建，在具体使用的时候才会创建。gorm内部，准确的说是database/sql内部会维护一个连接池，可以通过参数设置最大空闲连接数，连接最大空闲时间等。使用者不需要管连接的创建和关闭。
			err := sql.Ping() //看看是否被ping通
			if err != nil {
				log.SugarLogger.Error(err)
			} else {
				//设置连接数
				sql.SetMaxOpenConns(100)
				//设置连接池中的最大闲置连接数
				sql.SetMaxIdleConns(50)
				sql.SetConnMaxLifetime(30 * time.Second)
				time.Sleep(30 * time.Second)
			}
		}
	}()
}
