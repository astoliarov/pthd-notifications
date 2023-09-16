package rqueue

import (
	"errors"
	"github.com/rs/zerolog/log"
	"pthd-notifications/pkg/async-api/rqueue/messages"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/model"
)

//go:generate mockgen -destination=../../../tests/mocks/async-api/rqueue/iexecutor_mock.go -package=mocks pthd-notifications/pkg/async-api/rqueue IExecutor
type IExecutor interface {
	SendNotification(msg messages.RedisEventMessage) error
}

type SingleGoroutineExecutor struct {
	service *domain.Service
}

func NewSingleGoroutineExecutor(service *domain.Service) *SingleGoroutineExecutor {
	return &SingleGoroutineExecutor{
		service: service,
	}
}

func (executor *SingleGoroutineExecutor) SendNotification(msg messages.RedisEventMessage) error {
	notificationContext := msg.ToContext()

	executionErr := executor.service.SendNotification(notificationContext)
	var errNoSettings *domain.ErrNoSettings
	var errNoMessage *model.ErrNoMessage
	switch {
	case errors.As(executionErr, &errNoSettings):
		log.Info().Str("type", notificationContext.GetType()).Msg("No config for type")
		return nil
	case errors.As(executionErr, &errNoMessage):
		log.Info().Str("type", notificationContext.GetType()).Msg("No message for type")
		return nil
	}
	return executionErr
}
