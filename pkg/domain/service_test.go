package domain

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"pthd-notifications/pkg/domain/entities"
	"pthd-notifications/tests/mocks"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite

	controller     *gomock.Controller
	connectorMock  *mocks.MockIConnector
	repositoryMock *mocks.MockINotificationsRepo

	service *Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())

	s.connectorMock = mocks.NewMockIConnector(s.controller)
	s.repositoryMock = mocks.NewMockINotificationsRepo(s.controller)

	s.service = NewService(
		s.repositoryMock,
		s.connectorMock,
	)
}

func (s *ServiceTestSuite) Test_SendNotification_RepoErr_Err() {
	discordId := int64(2)
	err := errors.New("Test error")
	ctx := entities.NotificationContext{}

	s.repositoryMock.EXPECT().GetByDiscordId(discordId).Return(nil, err)

	sendErr := s.service.SendNotification(discordId, &ctx)

	assert.Equal(s.T(), sendErr, err)
}

func (s *ServiceTestSuite) Test_SendNotification_NoConfig_Err() {
	discordId := int64(2)
	ctx := entities.NotificationContext{}

	s.repositoryMock.EXPECT().GetByDiscordId(discordId).Return(nil, nil)

	sendErr := s.service.SendNotification(discordId, &ctx)

	assert.Nil(s.T(), sendErr)
}

func (s *ServiceTestSuite) Test_SendNotification_CreateNotificationErr_Err() {
	discordId := int64(2)
	ctx := entities.NotificationContext{}
	settings := &entities.NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined in chat",
		},
	}

	s.repositoryMock.EXPECT().GetByDiscordId(discordId).Return(settings, nil)

	sendErr := s.service.SendNotification(discordId, &ctx)

	assert.NotNil(s.T(), sendErr)
}

func (s *ServiceTestSuite) Test_SendNotification_SendErr_Err() {
	discordId := int64(2)
	ctx := entities.NotificationContext{}
	settings := &entities.NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined}} in chat",
		},
	}
	err := errors.New("Send error")

	s.repositoryMock.EXPECT().GetByDiscordId(discordId).Return(settings, nil)
	s.connectorMock.EXPECT().Send(gomock.Any()).Return(err)

	sendErr := s.service.SendNotification(discordId, &ctx)

	assert.Equal(s.T(), sendErr, err)
}

func (s *ServiceTestSuite) Test_SendNotification_NoErr_Nil() {
	discordId := int64(2)
	ctx := entities.NotificationContext{}
	settings := &entities.NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined}} in chat",
		},
	}

	s.repositoryMock.EXPECT().GetByDiscordId(discordId).Return(settings, nil)
	s.connectorMock.EXPECT().Send(gomock.Any()).Return(nil)

	sendErr := s.service.SendNotification(discordId, &ctx)

	assert.Nil(s.T(), sendErr)
}

func TestServiceTestSuite(t *testing.T) {
	testSuite := ServiceTestSuite{}
	suite.Run(t, &testSuite)
}
