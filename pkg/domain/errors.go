package domain

import (
	"fmt"
	"pthd-notifications/pkg/domain/model"
)

type ErrNoSettings struct {
	DiscordId int64
	Type      string
}

func (err *ErrNoSettings) Error() string {
	return fmt.Sprintf("no settings for DiscordId:%d and Type:%s", err.DiscordId, err.Type)
}

func NewErrNoSettingsFromContext(nc model.INotificationContext) error {
	return &ErrNoSettings{
		DiscordId: nc.GetDiscordId(),
		Type:      nc.GetType(),
	}
}
