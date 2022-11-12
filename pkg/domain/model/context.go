package model

const NotificationTypeUsersConnected = "users_connected"
const NotificationTypeUsersLeftChannel = "users_left_channel"

var NotificationTypes = []string{
	NotificationTypeUsersConnected,
	NotificationTypeUsersLeftChannel,
}

func IsNotificationTypeSupported(nType string) bool {
	for _, value := range NotificationTypes {
		if value == nType {
			return true
		}
	}

	return false
}

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
	return NotificationTypeUsersConnected
}

func (nc *UsersConnectedNotificationContext) GetDiscordId() int64 {
	return nc.Id
}

type UsersLeftChannelNotificationContext struct {
	Id int64
}

func (nc *UsersLeftChannelNotificationContext) GetType() string {
	return NotificationTypeUsersLeftChannel
}

func (nc *UsersLeftChannelNotificationContext) GetDiscordId() int64 {
	return nc.Id
}
