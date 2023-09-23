package messages

import (
	"pthd-notifications/pkg/domain/model"
	"strings"
	"time"
)

const MessageTypeNewUser = "new_user"
const MessageTypeUsersConnected = "users_connected"
const MessageTypeUsersLeft = "users_left"
const MessageTypeUserLeft = "user_left"

type MinimalMessage struct {
	MessageType string `json:"type"`
}

type RedisEventMessage interface {
	ToContext() model.INotificationContext
	GetHappenedAt() time.Time
}

type Message struct {
	Version     int       `json:"version"`
	MessageType string    `json:"type"`
	HappenedAt  time.Time `json:"happened_at"`
	ChannelId   int64     `json:"channel_id"`
}

func (msg Message) GetHappenedAt() time.Time {
	return msg.HappenedAt
}

type MessageNewUserInChannelData struct {
	Message
	Data struct {
		Username string `json:"username"`
	} `json:"data"`
}

func (msg MessageNewUserInChannelData) ToContext() model.INotificationContext {
	return &model.NewUserInChannelContext{
		Id:       msg.ChannelId,
		Username: msg.Data.Username,
	}
}

type MessageUsersConnectedToChannel struct {
	Message
	Data struct {
		Usernames []string `json:"usernames"`
	} `json:"data"`
}

func (msg MessageUsersConnectedToChannel) ToContext() model.INotificationContext {
	return &model.UsersConnectedNotificationContext{
		Id:          msg.ChannelId,
		Names:       msg.Data.Usernames,
		NamesJoined: strings.Join(msg.Data.Usernames, ","),
	}
}

type MessageUsersLeftChannel struct {
	Message
}

func (msg MessageUsersLeftChannel) ToContext() model.INotificationContext {
	return &model.UsersLeftChannelNotificationContext{
		Id: msg.ChannelId,
	}
}

type MessageUserLeftChannel struct {
	Message
	Data struct {
		Username string `json:"username"`
	} `json:"data"`
}

func (msg MessageUserLeftChannel) ToContext() model.INotificationContext {
	return &model.UserLeftChannelContext{
		Id:       msg.ChannelId,
		Username: msg.Data.Username,
	}
}
