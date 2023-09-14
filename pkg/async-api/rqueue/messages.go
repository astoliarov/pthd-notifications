package rqueue

import (
	"pthd-notifications/pkg/domain/model"
	"strings"
	"time"
)

const messageTypeNewUser = "new_user"
const messageTypeUsersConnected = "users_connected"
const messageTypeUsersLeave = "users_leave"

type minimalMessage struct {
	MessageType string `json:"type"`
}

type message struct {
	Version     int       `json:"version"`
	MessageType string    `json:"type"`
	HappenedAt  time.Time `json:"happened_at"`
	ChannelId   int64     `json:"channel_id"`
}

type iMessageWithContext interface {
	toContext() model.INotificationContext
}

type messageNewUserInChannelData struct {
	message
	Data struct {
		Username string `json:"username"`
	} `json:"data"`
}

func (msg messageNewUserInChannelData) toContext() model.INotificationContext {
	return &model.NewUserInChannelContext{
		Id:       msg.ChannelId,
		Username: msg.Data.Username,
	}
}

type messageUsersConnectedToChannel struct {
	message
	Data struct {
		Usernames []string `json:"usernames"`
	} `json:"data"`
}

func (msg messageUsersConnectedToChannel) toContext() model.INotificationContext {
	return &model.UsersConnectedNotificationContext{
		Id:          msg.ChannelId,
		Names:       msg.Data.Usernames,
		NamesJoined: strings.Join(msg.Data.Usernames, ","),
	}
}

type messageUsersLeftChannel struct {
	message
}

func (msg messageUsersLeftChannel) toContext() model.INotificationContext {
	return &model.UsersLeftChannelNotificationContext{
		Id: msg.ChannelId,
	}
}
