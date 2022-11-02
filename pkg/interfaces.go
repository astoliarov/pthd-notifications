package pkg

import (
	_ "github.com/golang/mock/mockgen/model"
	"pthd-notifications/pkg/entities"
)

//go:generate mockgen -destination=../tests/mocks/inotificationsrepo_mock.go -package=mocks pthd-notifications/pkg INotificationsRepo
type INotificationsRepo interface {
	GetByDiscordId(discordId int64) (*entities.NotificationSettings, error)
}

//go:generate mockgen -destination=../tests/mocks/iconnector_mock.go -package=mocks pthd-notifications/pkg IConnector
type IConnector interface {
	Send(notification *entities.Notification) error
}
