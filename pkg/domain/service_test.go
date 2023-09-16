package domain

import (
	"errors"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"pthd-notifications/pkg/domain/model"
	domain_mocks "pthd-notifications/tests/mocks/domain"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite

	controller     *gomock.Controller
	connectorMock  *domain_mocks.MockIConnector
	repositoryMock *domain_mocks.MockINotificationsRepo

	service *Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.controller = gomock.NewController(s.T())

	s.connectorMock = domain_mocks.NewMockIConnector(s.controller)
	s.repositoryMock = domain_mocks.NewMockINotificationsRepo(s.controller)

	s.service = NewService(
		s.repositoryMock,
		s.connectorMock,
	)
}

func (s *ServiceTestSuite) Test_SendNotification_RepoErr_Err() {
	discordId := int64(2)
	nType := model.NotificationTypeUsersConnected
	err := errors.New("test error")
	ctx := model.UsersConnectedNotificationContext{
		NamesJoined: "",
		Names:       []string{},
		Id:          discordId,
	}

	s.repositoryMock.EXPECT().Get(discordId, nType).Return(nil, err)

	sendErr := s.service.SendNotification(&ctx)

	assert.Equal(s.T(), sendErr, err)
}

func (s *ServiceTestSuite) Test_SendNotification_NoConfig_Err() {
	discordId := int64(2)
	ctx := model.UsersConnectedNotificationContext{
		NamesJoined: "",
		Names:       []string{},
		Id:          discordId,
	}
	nType := model.NotificationTypeUsersConnected

	s.repositoryMock.EXPECT().Get(discordId, nType).Return(nil, nil)

	sendErr := s.service.SendNotification(&ctx)

	assert.NotNil(s.T(), sendErr)
	assert.Equal(s.T(), &ErrNoSettings{
		DiscordId: discordId,
		Type:      ctx.GetType(),
	}, sendErr)
}

func (s *ServiceTestSuite) Test_SendNotification_CreateNotificationErr_Err() {
	discordId := int64(2)
	ctx := model.UsersConnectedNotificationContext{
		NamesJoined: "",
		Names:       []string{},
		Id:          discordId,
	}
	settings := &model.NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined in chat",
		},
		Type: model.NotificationTypeUsersConnected,
	}
	nType := model.NotificationTypeUsersConnected

	s.repositoryMock.EXPECT().Get(discordId, nType).Return(settings, nil)

	sendErr := s.service.SendNotification(&ctx)

	assert.NotNil(s.T(), sendErr)
}

func (s *ServiceTestSuite) Test_SendNotification_SendErr_Err() {
	discordId := int64(2)
	ctx := model.UsersConnectedNotificationContext{
		NamesJoined: "",
		Names:       []string{},
		Id:          discordId,
	}
	settings := &model.NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined}} in chat",
		},
		Type: model.NotificationTypeUsersConnected,
	}
	nType := model.NotificationTypeUsersConnected
	err := errors.New("test error")

	s.repositoryMock.EXPECT().Get(discordId, nType).Return(settings, nil)
	s.connectorMock.EXPECT().Send(gomock.Any()).Return(err)

	sendErr := s.service.SendNotification(&ctx)

	assert.Equal(s.T(), sendErr, err)
}

func (s *ServiceTestSuite) Test_SendNotification_NoErr_Nil() {
	discordId := int64(2)
	ctx := model.UsersConnectedNotificationContext{
		NamesJoined: "",
		Names:       []string{},
		Id:          discordId,
	}
	settings := &model.NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined}} in chat",
		},
		Type: model.NotificationTypeUsersConnected,
	}
	nType := model.NotificationTypeUsersConnected

	s.repositoryMock.EXPECT().Get(discordId, nType).Return(settings, nil)
	s.connectorMock.EXPECT().Send(gomock.Any()).Return(nil)

	sendErr := s.service.SendNotification(&ctx)

	assert.Nil(s.T(), sendErr)
}

func TestServiceTestSuite(t *testing.T) {
	testSuite := ServiceTestSuite{}
	suite.Run(t, &testSuite)
}
