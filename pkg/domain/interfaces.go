package domain

import (
	_ "go.uber.org/mock/mockgen/model"
	"pthd-notifications/pkg/domain/model"
)

//go:generate mockgen -destination=../../tests/mocks/domain/inotificationsrepo_mock.go -package=mocks pthd-notifications/pkg/domain INotificationsRepo
type INotificationsRepo interface {
	Get(discordId int64, notificationType string) (*model.NotificationSettings, error)
}

//go:generate mockgen -destination=../../tests/mocks/domain/iconnector_mock.go -package=mocks pthd-notifications/pkg/domain IConnector
type IConnector interface {
	Send(notification *model.Notification) error
}

//go:generate mockgen -destination=../../tests/mocks/domain/iservice_mock.go -package=mocks pthd-notifications/pkg/domain IService
type IService interface {
	SendNotification(notificationContext model.INotificationContext) error
}
