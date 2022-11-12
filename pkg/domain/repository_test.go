package domain

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pthd-notifications/pkg/domain/model"
	"testing"
)

type NotificationsMemoryRepositoryTestSuite struct {
	suite.Suite

	repo *NotificationsMemoryRepo
}

func (s *NotificationsMemoryRepositoryTestSuite) SetupSuite() {
	s.repo = NewNotificationsMemoryRepo()
}

func (s *NotificationsMemoryRepositoryTestSuite) Test__LoadSettings__GetReturnedCorrect() {
	settingsExample := &model.NotificationSettings{
		DiscordId:         1,
		TelegramChatId:    2,
		MessagesTemplates: []string{},
		Type:              model.UsersConnectedType,
	}
	settingsExample2 := &model.NotificationSettings{
		DiscordId:         2,
		TelegramChatId:    3,
		MessagesTemplates: []string{},
		Type:              model.UsersConnectedType,
	}

	settings := []*model.NotificationSettings{
		settingsExample,
		settingsExample2,
	}

	s.repo.Load(settings)

	returnedExample, err := s.repo.Get(2, model.UsersConnectedType)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), returnedExample.DiscordId, int64(2))
	assert.Equal(s.T(), returnedExample.TelegramChatId, int64(3))
	assert.Equal(s.T(), returnedExample.Type, model.UsersConnectedType)
}

func (s *NotificationsMemoryRepositoryTestSuite) Test__LoadSettings__GetTheSameIDButDifferentTypeReturnedCorrect() {
	settingsExample := &model.NotificationSettings{
		DiscordId:         1,
		TelegramChatId:    2,
		MessagesTemplates: []string{},
		Type:              model.UsersConnectedType,
	}
	settingsExample2 := &model.NotificationSettings{
		DiscordId:         1,
		TelegramChatId:    3,
		MessagesTemplates: []string{},
		Type:              model.UsersLeftChannelType,
	}

	settings := []*model.NotificationSettings{
		settingsExample,
		settingsExample2,
	}

	s.repo.Load(settings)

	returnedExample, err := s.repo.Get(1, model.UsersLeftChannelType)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), returnedExample.DiscordId, int64(1))
	assert.Equal(s.T(), returnedExample.TelegramChatId, int64(3))
	assert.Equal(s.T(), returnedExample.Type, model.UsersLeftChannelType)
}

func TestNotificationsMemoryRepo(t *testing.T) {
	testSuite := NotificationsMemoryRepositoryTestSuite{}
	suite.Run(t, &testSuite)
}
