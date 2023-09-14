package rqueue

import (
	"errors"
	"github.com/rs/zerolog/log"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/model"
)

type IExecutor interface {
	SendNotification(notificationContext model.INotificationContext) error
}

type SingleGoroutineExecutor struct {
	service *domain.Service
}

func NewSingleGoroutineExecutor(service *domain.Service) *SingleGoroutineExecutor {
	return &SingleGoroutineExecutor{
		service: service,
	}
}

func (executor *SingleGoroutineExecutor) SendNotification(notificationContext model.INotificationContext) error {
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
