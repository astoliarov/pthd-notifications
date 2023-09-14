package rqueue

import (
	"github.com/kelseyhightower/envconfig"
)

type RedisConfig struct {
	Host        string `default:"localhost" required:"false"`
	Port        int    `default:"6379" required:"false"`
	Db          int    `default:"0" required:"false"`
	Password    string `default:"" required:"false"`
	Queue       string `default:"notification.events" required:"true"`
	ReadTimeout int    `default:"1" required:"false"`
}

func LoadRedisConfig() (*RedisConfig, error) {
	var cfg RedisConfig
	err := envconfig.Process("notifications_redis", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
