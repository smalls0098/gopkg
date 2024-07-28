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
	return q.PushStr(ctx, key, string(data))
}

func (q *Queue) Pop(ctx context.Context, key string, v interface{}) error {
	res, err := q.PopStr(ctx, key)
	if len(res) == 0 {
		return nil
	}
	if err = json.Unmarshal([]byte(res), v); err != nil {
		return err
	}
	return nil
}

func (q *Queue) PushStr(ctx context.Context, key string, str string) error {
	_, err := q.rdb.LPush(ctx, key, str).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}
	return err
}

func (q *Queue) PopStr(ctx context.Context, key string) (string, error) {
	res, err := q.rdb.LPop(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}
	return res, nil
}
