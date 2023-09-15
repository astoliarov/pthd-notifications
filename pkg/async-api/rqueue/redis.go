package rqueue

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:generate mockgen -destination=../../../tests/mocks/iredisconnector_mock.go -package=mocks pthd-notifications/pkg/async-api/rqueue IRedisConnector
type IRedisConnector interface {
	ReadFromQueue(ctx context.Context) (string, string, error)
}

type RedisConnector struct {
	client    *redis.Client
	queueName string
	timeout   time.Duration
}

func NewRedisConnector(config *RedisConfig) *RedisConnector {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})

	return &RedisConnector{
		client:    client,
		queueName: config.Queue,
		timeout:   time.Second * time.Duration(config.ReadTimeout),
	}
}

func (connector *RedisConnector) ReadFromQueue(ctx context.Context) (string, string, error) {
	result := connector.client.BLPop(ctx, connector.timeout, connector.queueName)
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
