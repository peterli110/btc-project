package redisclient

import (
	"btc-project/library/log"
	"github.com/go-redis/redis"
)


type Config struct {
	Addr       string `toml:"addr"`
	Password   string `toml:"password"`
	DB         int    `toml:"db"`
	MaxRetries int    `toml:"maxretries"`
}

func NewRedis(c *Config) (client *redis.Client) {
	if c.Addr == "" {
		panic("redis address not defined")
	}

	client = redis.NewClient(&redis.Options{
		Addr: c.Addr,
		Password: c.Password,
		DB: c.DB,
		MaxRetries: c.MaxRetries,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Error("ping redis error: (%v)", err)
		panic(err)
	}

	log.Info("redis connected")
	return
}