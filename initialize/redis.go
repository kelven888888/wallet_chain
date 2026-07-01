package initialize

import (
	"context"
	"github.com/redis/go-redis/v9"

	"wallet_chain.com/global"

	"go.uber.org/zap"
)

func Redis() {
	redisCfg := global.SHOP_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.SHOP_LOG.Error("redis connect ping failed, err:", zap.Error(err))
	} else {
		global.SHOP_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		global.SHOP_REDIS = client
	}
	Walletinit()
}
