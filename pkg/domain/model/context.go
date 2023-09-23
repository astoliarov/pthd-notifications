package model

const NotificationTypeUsersConnected = "users_connected"
const NotificationTypeUsersLeftChannel = "users_left_channel"
const NotificationTypeNewUserInChannel = "new_user_in_channel"
const NotificationTypeNewUserLeftChannel = "user_left_channel"

var NotificationTypes = []string{
	NotificationTypeUsersConnected,
	NotificationTypeUsersLeftChannel,
	NotificationTypeNewUserInChannel,
	NotificationTypeNewUserLeftChannel,
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

type NewUserInChannelContext struct {
	Id       int64
	Username string
}

func (nc *NewUserInChannelContext) GetType() string {
	return NotificationTypeNewUserInChannel
}

func (nc *NewUserInChannelContext) GetDiscordId() int64 {
	return nc.Id
}

type UserLeftChannelContext struct {
	Id       int64
	Username string
}

func (nc *UserLeftChannelContext) GetType() string {
	return NotificationTypeNewUserLeftChannel
}

func (nc *UserLeftChannelContext) GetDiscordId() int64 {
	return nc.Id
}
