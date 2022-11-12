package domain

import (
	_ "github.com/golang/mock/mockgen/model"
	"pthd-notifications/pkg/domain/model"
)

//go:generate mockgen -destination=../../tests/mocks/inotificationsrepo_mock.go -package=mocks pthd-notifications/pkg/domain INotificationsRepo
type INotificationsRepo interface {
	Get(discordId int64, notificationType string) (*model.NotificationSettings, error)
}

//go:generate mockgen -destination=../../tests/mocks/iconnector_mock.go -package=mocks pthd-notifications/pkg/domain IConnector
type IConnector interface {
	Send(notification *model.Notification) error
}
