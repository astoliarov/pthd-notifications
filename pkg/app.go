package pkg

import (
	"context"
	"fmt"
	"pthd-notifications/pkg/api"
	"pthd-notifications/pkg/async-api/rqueue"
	"pthd-notifications/pkg/connectors"
	"pthd-notifications/pkg/domain"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

type Application struct {
	service *domain.Service
	config  *Config
}

func NewApplication() (*Application, error) {
	config, configErr := LoadConfig()
	if configErr != nil {
		return nil, fmt.Errorf("failed to load config: %s", configErr)
	}

	setupLogs()

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

	return &Application{
		service: service,
		config:  config,
	}, nil
}

func (app *Application) RunAPI(ctx context.Context) error {
	log.Info().Msg("Starting application")

	apiConfig, apiConfigErr := api.LoadApiConfig()
	if apiConfigErr != nil {
		return fmt.Errorf("failed to load API config: %s", apiConfigErr)
	}

	server := api.NewServer(apiConfig.Host, apiConfig.Port, app.config.Debug, app.service)

	return server.Run(ctx)
}

func (app *Application) RunRedisConsumer(ctx context.Context) error {
	log.Info().Msg("Starting Async application")
	redisConfig, redisConfigErr := rqueue.LoadRedisConfig()
	if redisConfigErr != nil {
		return fmt.Errorf("failed to load redis config: %s", redisConfigErr)
	}

	asyncExecutor := rqueue.NewSingleGoroutineExecutor(app.service)
	asyncApi := rqueue.NewRedisAsyncAPI(asyncExecutor, redisConfig)

	return asyncApi.RunConsumer(ctx)
}
