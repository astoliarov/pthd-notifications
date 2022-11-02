package connectors

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"pthd-notifications/pkg/entities"
)

func InitBot(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return bot, nil
}

type TelegramConnector struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramConnector(bot *tgbotapi.BotAPI) *TelegramConnector {
	return &TelegramConnector{bot: bot}
}

func (tc *TelegramConnector) Send(notification *entities.Notification) error {
	msg := tgbotapi.NewMessage(notification.TelegramChatId, notification.Message)
	_, sendErr := tc.bot.Send(msg)
	return sendErr
}
