package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken string `required:"true"`
	ApiPort       int    `default:"3030" required:"true"`
	ApiHost       string `default:"0.0.0.0" required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("notifications", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
