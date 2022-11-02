package domain

import (
	"pthd-notifications/pkg/domain/entities"
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

func (s *Service) SendNotification(discordId int64, notificationContext *entities.NotificationContext) error {
	settings, repoErr := s.repo.GetByDiscordId(discordId)
	if repoErr != nil {
		return repoErr
	}
	// TODO: maybe here it is better to throw an error?
	if settings == nil {
		return nil
	}

	notification, createErr := entities.NewNotification(notificationContext, settings)
	if createErr != nil {
		return createErr
	}
	// TODO: maybe here it is better to throw an error?
	if notification == nil {
		return nil
	}

	sendErr := s.connector.Send(notification)
	if sendErr != nil {
		return sendErr
	}

	return nil
}
