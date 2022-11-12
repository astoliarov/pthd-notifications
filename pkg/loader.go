package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"pthd-notifications/pkg/domain/model"
)

type NotificationSettingsModel struct {
	DiscordId         int64    `json:"discord_id"`
	TelegramChatId    int64    `json:"telegram_chat_id"`
	MessagesTemplates []string `json:"messages_templates"`
	Type              string   `json:"type"`
}

func (ns *NotificationSettingsModel) ToEntity() (*model.NotificationSettings, error) {
	if !model.IsNotificationTypeSupported(ns.Type) {
		return nil, errors.New(
			fmt.Sprintf("Notification with type:%s and discord_id: %d not supported", ns.Type, ns.DiscordId),
		)
	}

	return &model.NotificationSettings{
		DiscordId:         ns.DiscordId,
		TelegramChatId:    ns.TelegramChatId,
		MessagesTemplates: ns.MessagesTemplates,
		Type:              ns.Type,
	}, nil
}

type NotificationSettingsFile struct {
	Settings []NotificationSettingsModel `json:"settings"`
}

type Loader struct {
	path string
}

func NewLoader(path string) *Loader {
	return &Loader{
		path: path,
	}
}

func (l *Loader) Load() ([]*model.NotificationSettings, error) {
	dat, readErr := os.ReadFile(l.path)
	if readErr != nil {
		return nil, readErr
	}

	settingsFile := NotificationSettingsFile{}

	unmarshalErr := json.Unmarshal(dat, &settingsFile)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	var ents []*model.NotificationSettings

	for _, mdl := range settingsFile.Settings {
		entity, toEntErr := mdl.ToEntity()
		if toEntErr != nil {
			log.Println(toEntErr.Error())
			continue
		}

		ents = append(ents, entity)
	}
	return ents, nil
}
