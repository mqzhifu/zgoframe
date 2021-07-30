package initialize

import (
	"go.uber.org/zap"
	"zgoframe/core/global"
	"github.com/go-redis/redis"
)

func GetNewRedis()(*redis.Client,error) {
	redisCfg := global.C.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Ip + ":"+ redisCfg.Port,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DbNumber,       // use default DB
	})
	pong, err := client.Ping().Result()

	if err != nil {
		global.V.Zap.Error("redis connect ping failed, err:", zap.Any("err", err))
		return client,err
	}

	global.V.Zap.Info("redis connect ping response:", zap.String("pong",pong))
	return client,nil
}

