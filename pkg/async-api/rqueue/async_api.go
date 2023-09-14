package rqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RedisAsyncAPI struct {
	client    *redis.Client
	queueName string
	timeout   time.Duration
	executor  IExecutor
}

func NewRedisAsyncAPI(executor IExecutor, config *RedisConfig) *RedisAsyncAPI {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})

	return &RedisAsyncAPI{
		client:    client,
		queueName: config.Queue,
		timeout:   time.Second * time.Duration(config.ReadTimeout),
		executor:  executor,
	}
}

func (asyncApi *RedisAsyncAPI) RunConsumer(ctx context.Context) error {

	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	isRunning := true
	for isRunning {
		select {
		case <-ctx.Done():
			isRunning = false
			log.Debug().Msg("stop consumer by context.Done()")
		case <-signalCh:
			isRunning = false
			log.Debug().Msg("stop consumer by signal")
		default:
			asyncApi.executeRead(ctx)
		}
	}

	return nil
}

func (asyncApi *RedisAsyncAPI) executeRead(ctx context.Context) {
	log.Debug().Msg("reading data from queue")
	_, data, resultErr := asyncApi.readFromQueue(ctx)
	if resultErr != nil {
		log.Error().Err(resultErr).Msg("ReadFromQueue error")
		sentry.CaptureException(resultErr)
		return
	}

	processErr := asyncApi.processMessage(data)
	if processErr != nil {
		log.Error().Err(processErr).Msg("processMessage error")
		sentry.CaptureException(resultErr)
	}
}

func (asyncApi *RedisAsyncAPI) readFromQueue(ctx context.Context) (string, string, error) {
	result := asyncApi.client.BLPop(ctx, asyncApi.timeout, asyncApi.queueName)
	msg, err := result.Result()
	if err != nil && err != redis.Nil {
		return "", "", err
	}

	if len(msg) == 0 {
		return "", "", nil
	}

	if len(msg) != 2 {
		return "", "", ErrUnexpectedResponseLength
	}

	return msg[0], msg[1], nil
}

func (asyncApi *RedisAsyncAPI) processMessage(data string) error {
	if data == "" {
		return nil
	}

	minMessage := &minimalMessage{}
	unmarshalErr := json.Unmarshal([]byte(data), minMessage)
	if unmarshalErr != nil {
		return unmarshalErr
	}

	var message iMessageWithContext

	switch minMessage.MessageType {
	case messageTypeNewUser:
		message = &messageNewUserInChannelData{}
	case messageTypeUsersConnected:
		message = &messageUsersConnectedToChannel{}
	case messageTypeUsersLeave:
		message = &messageUsersLeftChannel{}
	default:
		return ErrUnexpectedMessageType
	}

	unmarshalFullMessageErr := json.Unmarshal([]byte(data), message)
	if unmarshalFullMessageErr != nil {
		return unmarshalFullMessageErr
	}

	return asyncApi.executor.SendNotification(message.toContext())
}
