package initialization

import (
	"srv/log"
	"srv/models"

	"github.com/gomodule/redigo/redis"
)

var Redis *redis.Pool

func InitRedis() { //redis连接池
	Redis = &redis.Pool{ //现在先用着，不会的后边查
		MaxIdle:     100,
		MaxActive:   0,
		IdleTimeout: 100,
		Wait:        true, //超过最大连接等待
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", models.Overall_Situation_Redis.Redis.Port)
			if err != nil {
				log.SugarLogger.Panic(err)
			}
			if _, err := c.Do("AUTH", models.Overall_Situation_Redis.Redis.Password); err != nil {
				c.Close()
				log.SugarLogger.Panic(err)
			}
			return c, nil
		},
	}
}
