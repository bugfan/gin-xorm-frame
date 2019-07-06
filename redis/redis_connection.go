package redis

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/redis.v5"
)

//  连接R
var R *Redis

func init() {
	R = new(Redis)
	R.RedisCluster = new(redis.Client)
}

func GetRedis() *Redis {
	return R
}

type Config struct {
	Addr       string
	PoolSize   int
	DB         int
	Password   string
	MaxRetries int
}

// Redis redis连接对象
type Redis struct {
	RedisCluster *redis.Client
}

// 链接redis
func SetEngine(c *Config) (int, error) {
	if c.MaxRetries == 0 {
		c.MaxRetries = 16
	}
	addr := c.Addr
	if addrs := strings.Split(addr, ":"); len(addrs) < 2 {
		addr = addrs[0] + ":6379"
	}
	options := redis.Options{
		Addr:        addr,
		MaxRetries:  c.MaxRetries, //最大重试次数，默认16
		PoolSize:    c.PoolSize,
		DB:          c.DB,
		ReadTimeout: 500 * time.Millisecond,
		IdleTimeout: 12 * time.Second,
	}
	if strings.TrimSpace(c.Password) != "" {
		options.Password = c.Password
	}
	R.RedisCluster = redis.NewClient(&options)
	if R.Set("Redis OK", 1, 1000000000) {
		log.Printf("Redis Init Succeeded! Index:[%v]\n", c.DB)
		return 0, nil
	} else {
		return -1, errors.New(fmt.Sprintf("Redis Init Failed! Datil:%v", options))
	}
}
