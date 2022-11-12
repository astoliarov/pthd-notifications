package model

import "fmt"

type ErrNoMessage struct {
	DiscordId int64
	Type      string
}

func (err *ErrNoMessage) Error() string {
	return fmt.Sprintf("no message template for DiscordId:%d and Type:%s", err.DiscordId, err.Type)
}

func NewErrNoSettingsFromContext(nc INotificationContext) error {
	return &ErrNoMessage{
		DiscordId: nc.GetDiscordId(),
		Type:      nc.GetType(),
	}
}
