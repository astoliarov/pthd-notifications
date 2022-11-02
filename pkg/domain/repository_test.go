package domain

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pthd-notifications/pkg/domain/entities"
	"testing"
)

type NotificationsMemoryRepositoryTestSuite struct {
	suite.Suite

	repo *NotificationsMemoryRepo
}

func (s *NotificationsMemoryRepositoryTestSuite) SetupSuite() {
	s.repo = NewNotificationsMemoryRepo()
}

func (s *NotificationsMemoryRepositoryTestSuite) Test__LoadSettings__GetByDiscordIdReturnedCorrect() {
	settingsExample := &entities.NotificationSettings{
		DiscordId:         1,
		TelegramChatId:    2,
		MessagesTemplates: []string{},
	}
	settingsExample2 := &entities.NotificationSettings{
		DiscordId:         2,
		TelegramChatId:    3,
		MessagesTemplates: []string{},
	}

	settings := []*entities.NotificationSettings{
		settingsExample,
		settingsExample2,
	}

	s.repo.Load(settings)

	returnedExample, err := s.repo.GetByDiscordId(2)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), returnedExample.DiscordId, int64(2))
	assert.Equal(s.T(), returnedExample.TelegramChatId, int64(3))
}

func TestNotificationsMemoryRepo(t *testing.T) {
	testSuite := NotificationsMemoryRepositoryTestSuite{}
	suite.Run(t, &testSuite)
}
