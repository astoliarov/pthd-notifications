package api

import "github.com/kelseyhightower/envconfig"

type ApiConfig struct {
	Port int    `default:"3030" required:"true"`
	Host string `default:"0.0.0.0" required:"true"`
}

func LoadApiConfig() (*ApiConfig, error) {
	var cfg ApiConfig
	err := envconfig.Process("notifications_api", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
