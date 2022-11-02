package pkg

import (
	"context"
	"fmt"
	"log"
	"pthd-notifications/pkg/connectors"
)

type Application struct {
	service *Service
}

func NewApplication() (*Application, error) {
	config, configErr := LoadConfig()
	if configErr != nil {
		return nil, fmt.Errorf("failed to load config: %s", configErr)
	}

	tgBot, botInitErr := connectors.InitBot(config.TelegramToken)
	if botInitErr != nil {
		return nil, fmt.Errorf("cannot initialize telegram bot: %s", botInitErr)
	}
	tgConnector := connectors.NewTelegramConnector(tgBot)
	repository := NewNotificationsMemoryRepo()

	service := NewService(repository, tgConnector)

	return &Application{
		service: service,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	log.Println("Running application")
	return nil
}
