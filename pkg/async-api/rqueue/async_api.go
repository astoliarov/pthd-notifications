package rqueue

import (
	"context"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
	"os/signal"
	"pthd-notifications/pkg/async-api/rqueue/messages"
	"syscall"
)

type RedisAsyncAPI struct {
	connector IRedisConnector
	executor  IExecutor
}

func NewRedisAsyncAPI(executor IExecutor, connector IRedisConnector) *RedisAsyncAPI {

	return &RedisAsyncAPI{
		connector: connector,
		executor:  executor,
	}
}

func (asyncApi *RedisAsyncAPI) RunConsumer(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	isRunning := true
	for isRunning {
		select {
		case <-ctx.Done():
			isRunning = false
			log.Debug().Msg("stop consumer by context.Done()")
		default:
			readErr := asyncApi.executeRead(ctx)
			// critical error happened, stop consumer
			if readErr != nil {
				return readErr
			}
		}
	}

	return nil
}

func (asyncApi *RedisAsyncAPI) executeRead(ctx context.Context) error {
	log.Debug().Msg("reading data from queue")
	_, data, resultErr := asyncApi.connector.ReadFromQueue(ctx)
	if resultErr != nil {
		log.Error().Err(resultErr).Msg("ReadFromQueue error")
		sentry.CaptureException(resultErr)
		return resultErr
	}

	processErr := asyncApi.processMessage(data)
	if processErr != nil {
		log.Error().Err(processErr).Msg("processMessage error")
		sentry.CaptureException(resultErr)
	}

	return nil
}

func (asyncApi *RedisAsyncAPI) processMessage(data string) error {
	if data == "" {
		return nil
	}

	minMessage := &messages.MinimalMessage{}
	unmarshalErr := json.Unmarshal([]byte(data), minMessage)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	var message messages.RedisEventMessage

	switch minMessage.MessageType {
	case messages.MessageTypeNewUser:
		message = &messages.MessageNewUserInChannelData{}
	case messages.MessageTypeUsersConnected:
		message = &messages.MessageUsersConnectedToChannel{}
	case messages.MessageTypeUsersLeft:
		message = &messages.MessageUsersLeftChannel{}
	case messages.MessageTypeUserLeft:
		message = &messages.MessageUserLeftChannel{}
	default:
		return ErrUnexpectedMessageType
	}

	unmarshalFullMessageErr := json.Unmarshal([]byte(data), message)
	if unmarshalFullMessageErr != nil {
		return unmarshalFullMessageErr
	}

	return asyncApi.executor.SendNotification(message)
}
