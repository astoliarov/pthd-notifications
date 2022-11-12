package domain

import (
	"fmt"
	"pthd-notifications/pkg/domain/model"
)

type NotificationsMemoryRepo struct {
	settingsByDiscordID map[string]*model.NotificationSettings
}

func NewNotificationsMemoryRepo() *NotificationsMemoryRepo {

	return &NotificationsMemoryRepo{
		settingsByDiscordID: map[string]*model.NotificationSettings{},
	}
}

func (repo *NotificationsMemoryRepo) getKey(discordId int64, notificationType string) string {
	return fmt.Sprintf("%d_%s", discordId, notificationType)
}

func (repo *NotificationsMemoryRepo) Get(discordId int64, notificationType string) (*model.NotificationSettings, error) {
	settings, ok := repo.settingsByDiscordID[repo.getKey(discordId, notificationType)]
	if !ok {
		return nil, nil
	}

	return settings, nil
}

func (repo *NotificationsMemoryRepo) Load(configs []*model.NotificationSettings) {
	for _, config := range configs {
		repo.settingsByDiscordID[repo.getKey(config.DiscordId, config.Type)] = config
	}
}
