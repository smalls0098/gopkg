package redis

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
)

type Queue struct {
	rdb *redis.Client
}

func NewQueue(rdb *redis.Client) *Queue {
	return &Queue{
		rdb: rdb,
	}
}

func (q *Queue) Push(ctx context.Context, key string, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = q.rdb.LPush(ctx, key, string(data)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return err
}

func (q *Queue) Pop(ctx context.Context, key string, v interface{}) error {
	res, err := q.rdb.LPop(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	if len(res) == 0 {
		return nil
	}
	if err = json.Unmarshal([]byte(res), v); err != nil {
		return err
	}
	return nil
}
