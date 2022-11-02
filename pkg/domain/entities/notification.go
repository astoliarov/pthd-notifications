package entities

import (
	"bytes"
	"math/rand"
	"text/template"
)

type NotificationContext struct {
	NamesJoined string
	Names       []string
}

type NotificationSettings struct {
	DiscordId         int64
	TelegramChatId    int64
	MessagesTemplates []string
}

func (ns *NotificationSettings) getRandomMessageTemplate() *string {
	if len(ns.MessagesTemplates) == 0 {
		return nil
	}

	tmpl := ns.MessagesTemplates[rand.Intn(len(ns.MessagesTemplates))]
	return &tmpl
}

type Notification struct {
	TelegramChatId int64
	Message        string
}

func NewNotification(notificationContext *NotificationContext, settings *NotificationSettings) (*Notification, error) {
	messageTemplate := settings.getRandomMessageTemplate()
	if messageTemplate == nil {
		return nil, nil
	}

	tmpl, parseErr := template.New("").Parse(*messageTemplate)
	if parseErr != nil {
		return nil, parseErr
	}

	var tpl bytes.Buffer
	executeErr := tmpl.Execute(&tpl, notificationContext)
	if executeErr != nil {
		return nil, executeErr
	}

	return &Notification{
		TelegramChatId: settings.TelegramChatId,
		Message:        tpl.String(),
	}, nil
}
