package redis

import (
	"fmt"

	"100.GO/internal/entity/user"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{Client: rdb}
}

func (r *RedisClient) SetHash(key string, values map[string]interface{}) error {
	return r.Client.HMSet(key, values).Err()
}

func (r *RedisClient) VerifyEmail(email string, usercode int64) (*user.CreateUser, error) {
	code, err := r.Client.HGet(email, "code").Int64()
	if err != nil {
		fmt.Println(code, err)
		return nil, fmt.Errorf("error HGET:%v", err)
	}
	if code == usercode {
		result, err := r.Client.HGetAll(email).Result()
		if err != nil {
			return nil, fmt.Errorf("error HGETALL: %v", err)
		}
		return &user.CreateUser{
			Firstname: result["firstname"],
			Lastname:  result["lastname"],
			Email:     result["email"],
			Password:  result["password"],
		}, nil

	}
	return nil, err
}

func (r *RedisClient) GetHash(key string) (map[string]string, error) {
	return r.Client.HGetAll(key).Result()
}

func (r *RedisClient) Delete(key string) error {
	return r.Client.Del(key).Err()
}
