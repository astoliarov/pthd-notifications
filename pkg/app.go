package pkg

import (
	"context"
	"fmt"
	"log"
	"pthd-notifications/pkg/api"
	"pthd-notifications/pkg/connectors"
	"pthd-notifications/pkg/domain"

	"github.com/getsentry/sentry-go"
)

type Application struct {
	service *domain.Service
	server  *api.Server
}

func NewApplication() (*Application, error) {
	config, configErr := LoadConfig()
	if configErr != nil {
		return nil, fmt.Errorf("failed to load config: %s", configErr)
	}

	if config.SentryDSN != "" {
		if sentryInitErr := sentry.Init(sentry.ClientOptions{
			Dsn: config.SentryDSN,
		}); sentryInitErr != nil {
			return nil, fmt.Errorf("failed to initialize sentry: %s", configErr)
		}
	}

	settingsLoader := NewLoader(config.PathToSettings)
	settings, loadErr := settingsLoader.Load()
	if loadErr != nil {
		return nil, fmt.Errorf("failed to load settings: %s", loadErr)
	}

	tgBot, botInitErr := connectors.InitBot(config.TelegramToken, config.Debug)
	if botInitErr != nil {
		return nil, fmt.Errorf("cannot initialize telegram bot: %s", botInitErr)
	}
	tgConnector := connectors.NewTelegramConnector(tgBot)
	repository := domain.NewNotificationsMemoryRepo()

	repository.Load(settings)

	service := domain.NewService(repository, tgConnector)
	server := api.NewServer(config.ApiHost, config.ApiPort, service)

	return &Application{
		service: service,
		server:  server,
	}, nil
}

func (app *Application) Run(ctx context.Context) error {
	log.Println("Starting application")
	app.server.Run(ctx)
	return nil
}
