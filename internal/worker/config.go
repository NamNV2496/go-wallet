package worker

import (
	"github.com/hibiken/asynq"
	"github.com/namnv2496/go-wallet/config"
)

type RedisConfigOpt struct {
	redisOpt *asynq.RedisClientOpt
}

func NewRedisConfigOpt(
	config config.Config,
) *RedisConfigOpt {
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	return &RedisConfigOpt{
		redisOpt: &redisOpt,
	}
}
