package pkg

import "pthd-notifications/pkg/entities"

type NotificationsMemoryRepo struct {
	settingsByDiscordID map[int64]*entities.NotificationSettings
}

func NewNotificationsMemoryRepo() *NotificationsMemoryRepo {

	return &NotificationsMemoryRepo{
		settingsByDiscordID: map[int64]*entities.NotificationSettings{},
	}
}

func (repo *NotificationsMemoryRepo) GetByDiscordId(discordId int64) (*entities.NotificationSettings, error) {
	settings, ok := repo.settingsByDiscordID[discordId]
	if !ok {
		return nil, nil
	}

	return settings, nil
}

func (repo *NotificationsMemoryRepo) Load(configs []*entities.NotificationSettings) {
	for _, config := range configs {
		repo.settingsByDiscordID[config.DiscordId] = config
	}
}
