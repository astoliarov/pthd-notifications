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
	service  *domain.Service
	server   *api.Server
	asyncApi *rqueue.RedisAsyncAPI
}

func NewApplication() (*Application, error) {
	config, configErr := LoadConfig()
	if configErr != nil {
		return nil, fmt.Errorf("failed to load config: %s", configErr)
	}
	redisConfig, redisConfigErr := rqueue.LoadRedisConfig()
	if redisConfigErr != nil {
		return nil, fmt.Errorf("failed to load config: %s", redisConfigErr)
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
	server := api.NewServer(config.ApiHost, config.ApiPort, config.Debug, service)

	asyncExecutor := rqueue.NewSingleGoroutineExecutor(service)
	asyncApi := rqueue.NewRedisAsyncAPI(asyncExecutor, redisConfig)

	return &Application{
		service:  service,
		server:   server,
		asyncApi: asyncApi,
	}, nil
}

func (app *Application) RunAPI(ctx context.Context) error {
	log.Info().Msg("Starting application")
	return app.server.Run(ctx)
}

func (app *Application) RunRedisConsumer(ctx context.Context) error {
	log.Info().Msg("Starting Async application")
	return app.asyncApi.RunConsumer(ctx)
}
