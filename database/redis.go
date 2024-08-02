package database

import (
	"chanel/config"
	"chanel/lib"
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	config *config.Config
	ctx    context.Context
	tools  *lib.Tools

	Client *redis.Client
}

const (
	Maintain = "Maintain"
)

func RedisInit(config *config.Config, ctx context.Context, tools *lib.Tools) *Redis {
	defer func() {
		if err := recover(); err != nil {
			panic(fmt.Errorf("redis error -> %v", lib.PanicParser(err)))
		}
	}()

	return &Redis{
		config: config,
		ctx:    ctx,
		tools:  tools,
	}
}

func (r *Redis) Start() {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     r.config.RedisHost + ":" + r.config.RedisPort,
		DB:       0,
		PoolSize: r.tools.StrToInt(r.config.RedisPoolSize),
	})
	// 測試 PingPong
	_, err := r.Client.Ping(r.ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("redis start 錯誤, ERR: %s", err.Error()))
	}
}
