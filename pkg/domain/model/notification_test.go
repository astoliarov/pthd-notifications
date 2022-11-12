package model

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type NewNotificationTestSuite struct {
	suite.Suite
}

func (s *NewNotificationTestSuite) Test__AllCorrect__NotificationCreated() {
	settings := &NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined}} in chat",
		},
	}
	context := &UsersConnectedNotificationContext{
		Names:       []string{"КурочкаХлеба", "TiConcerto"},
		NamesJoined: "КурочкаХлеба,TiConcerto",
		Id:          2,
	}

	notification, err := NewNotification(
		context,
		settings,
	)

	assert.Equal(s.T(), "Let's battle. Already КурочкаХлеба,TiConcerto in chat", notification.Message)
	assert.Nil(s.T(), err)
}

func (s *NewNotificationTestSuite) Test__NoMessageInSettings__NoNotification() {
	settings := &NotificationSettings{
		DiscordId:         2,
		TelegramChatId:    2,
		MessagesTemplates: []string{},
		Type:              NotificationTypeUsersConnected,
	}
	context := &UsersConnectedNotificationContext{
		Names:       []string{"КурочкаХлеба", "TiConcerto"},
		NamesJoined: "КурочкаХлеба,TiConcerto",
		Id:          2,
	}

	notification, err := NewNotification(
		context,
		settings,
	)

	assert.Equal(s.T(), &ErrNoMessage{
		DiscordId: 2,
		Type:      NotificationTypeUsersConnected,
	}, err)
	assert.Nil(s.T(), notification)
}

func (s *NewNotificationTestSuite) Test__BrokenTemplate__Err() {
	settings := &NotificationSettings{
		DiscordId:      1,
		TelegramChatId: 2,
		MessagesTemplates: []string{
			"Let's battle. Already {{.NamesJoined in chat",
		},
	}
	context := &UsersConnectedNotificationContext{
		Names:       []string{"КурочкаХлеба", "TiConcerto"},
		NamesJoined: "КурочкаХлеба,TiConcerto",
		Id:          2,
	}

	notification, err := NewNotification(
		context,
		settings,
	)

	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), notification)
}

func TestNewNotification(t *testing.T) {
	testSuite := NewNotificationTestSuite{}
	suite.Run(t, &testSuite)
}
