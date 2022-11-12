package model

const UsersConnectedType = "users_connected"
const UsersLeftChannelType = "users_left_channel"

type INotificationContext interface {
	GetType() string
	GetDiscordId() int64
}

type UsersConnectedNotificationContext struct {
	NamesJoined string
	Names       []string
	Id          int64
}

func (nc *UsersConnectedNotificationContext) GetType() string {
	return UsersConnectedType
}

func (nc *UsersConnectedNotificationContext) GetDiscordId() int64 {
	return nc.Id
}

type UsersLeftChannelNotificationContext struct {
	Id int64
}

func (nc *UsersLeftChannelNotificationContext) GetType() string {
	return UsersLeftChannelType
}

func (nc *UsersLeftChannelNotificationContext) GetDiscordId() int64 {
	return nc.Id
}
