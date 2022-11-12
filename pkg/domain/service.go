package domain

import (
	"pthd-notifications/pkg/domain/model"
)

type Service struct {
	repo      INotificationsRepo
	connector IConnector
}

func NewService(repo INotificationsRepo, connector IConnector) *Service {
	return &Service{
		repo:      repo,
		connector: connector,
	}
}

func (s *Service) SendNotification(notificationContext model.INotificationContext) error {
	settings, repoErr := s.repo.Get(notificationContext.GetDiscordId(), notificationContext.GetType())
	if repoErr != nil {
		return repoErr
	}

	if settings == nil {
		return NewErrNoSettingsFromContext(notificationContext)
	}

	notification, createErr := model.NewNotification(notificationContext, settings)
	if createErr != nil {
		return createErr
	}
	if notification == nil {
		return nil
	}

	sendErr := s.connector.Send(notification)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
