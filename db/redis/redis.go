package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	stdLog "log"
	"time"
)

type Options struct {
	Addr         string
	Username     string
	Password     string
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func NewWithConf(o Options) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         o.Addr,
		Username:     o.Username,
		Password:     o.Password,
		DB:           o.DB,
		DialTimeout:  o.DialTimeout,
		WriteTimeout: o.WriteTimeout,
		ReadTimeout:  o.ReadTimeout,
	})

	stdLog.Printf("opening connection to redis")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("redis error: %s", err.Error()))
	} else {
		stdLog.Printf("redis client ping result: %s", p)
	}
	return rdb
}
