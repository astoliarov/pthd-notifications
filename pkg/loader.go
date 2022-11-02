package pkg

import (
	"encoding/json"
	"os"
	"pthd-notifications/pkg/domain/entities"
)

type NotificationSettingsModel struct {
	DiscordId         int64    `json:"discord_id"`
	TelegramChatId    int64    `json:"telegram_chat_id"`
	MessagesTemplates []string `json:"messages_templates"`
}

func (ns *NotificationSettingsModel) ToEntity() *entities.NotificationSettings {
	return &entities.NotificationSettings{
		DiscordId:         ns.DiscordId,
		TelegramChatId:    ns.TelegramChatId,
		MessagesTemplates: ns.MessagesTemplates,
	}
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

func (l *Loader) Load() ([]*entities.NotificationSettings, error) {
	dat, readErr := os.ReadFile(l.path)
	if readErr != nil {
		return nil, readErr
	}

	settingsFile := NotificationSettingsFile{}

	unmarshalErr := json.Unmarshal(dat, &settingsFile)
	if unmarshalErr != nil {
		return nil, unmarshalErr
	}

	var ents []*entities.NotificationSettings

	for _, model := range settingsFile.Settings {
		ents = append(ents, model.ToEntity())
	}

	return ents, nil
}
