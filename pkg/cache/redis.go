package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

var (
	redisClient *redis.Client
	once        sync.Once
)

func InitRedisClient() {

	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%s:%s", viper.GetString("redis.hostname"), viper.GetString("redis.port")),
			Password:     viper.GetString("redis.password"),
			DB:           viper.GetInt("redis.db"),
			DialTimeout:  time.Duration(viper.GetInt64("redis.dialTimeOut")) * time.Second,
			ReadTimeout:  time.Duration(viper.GetInt64("redis.readTimeOut")) * time.Second,
			WriteTimeout: time.Duration(viper.GetInt64("redis.writeTimeOut")) * time.Second,
			PoolSize:     viper.GetInt("redis.poolSize"),
			PoolTimeout:  time.Duration(viper.GetInt64("redis.poolTimeOut")) * time.Second,
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			log.Fatal("[error] redis ping error : ", err)
		}

		fmt.Println("[info] redis successfully initialized")
		redisClient = client
	})
}
func NewRedisClient() *redis.Client {
	return redisClient
}
func CloseRedisClient() {
	_ = redisClient.Close()
}

func Set(ctx context.Context, key, val string, ttl time.Duration) error {
	if err := redisClient.Set(ctx, key, val, ttl).Err(); err != nil {
		fmt.Println("[error] redis set key failed: ", err)
		return err
	}
	return nil
}

func Get(ctx context.Context, key string) string {
	result, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		fmt.Printf("[error] Key %s not found in Redis cache", key)
		return ""
	}
	return result
}

func Delete(ctx context.Context, key string) {
	redisClient.Del(ctx, key)
}
