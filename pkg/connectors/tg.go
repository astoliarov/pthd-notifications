package connectors

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"pthd-notifications/pkg/domain/model"
	"time"
)

func InitBot(token string, debug bool) (*tgbotapi.BotAPI, error) {
	client := http.Client{Timeout: 5 * time.Second}
	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, &client)
	if err != nil {
		return nil, err
	}

	if debug {
		bot.Debug = true
	}

	log.Info().Str("username", bot.Self.UserName).Msg("Authorized on account")

	return bot, nil
}

type TelegramConnector struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramConnector(bot *tgbotapi.BotAPI) *TelegramConnector {
	return &TelegramConnector{bot: bot}
}

func (tc *TelegramConnector) Send(notification *model.Notification) error {
	msg := tgbotapi.NewMessage(notification.TelegramChatId, notification.Message)
	_, sendErr := tc.bot.Send(msg)
	return sendErr
}
