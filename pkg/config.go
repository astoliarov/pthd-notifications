package pkg

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TelegramToken  string `required:"true"`
	PathToSettings string `default:"./settings.json"`
	Debug          bool   `default:"false"`
	SentryDSN      string
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("notifications", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
