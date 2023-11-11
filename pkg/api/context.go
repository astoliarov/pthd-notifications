package api

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"net/http"
	"pthd-notifications/pkg/domain/model"
	"strings"
)

type notificationTypeParams struct {
	Type string `schema:"type" validate:"required" json:"type"`
}

type usersConnectedNotificationContext struct {
	DiscordId int64  `schema:"discord_id" validate:"required" json:"discord_id"`
	Usernames string `schema:"usernames" validate:"required" json:"usernames"`
	Type      string `schema:"type" validate:"required" json:"type"`
}

func (np *usersConnectedNotificationContext) toContext() model.INotificationContext {
	names := strings.Split(np.Usernames, ",")

	return &model.UsersConnectedNotificationContext{
		NamesJoined: np.Usernames,
		Names:       names,
		Id:          np.DiscordId,
	}
}

type usersLeftChannelNotificationContext struct {
	DiscordId int64  `schema:"discord_id" validate:"required" json:"discord_id"`
	Type      string `schema:"type" validate:"required" json:"type"`
}

func (np *usersLeftChannelNotificationContext) toContext() model.INotificationContext {

	return &model.UsersLeftChannelNotificationContext{
		Id: np.DiscordId,
	}
}

type notificationContextParser struct {
	decoder   *schema.Decoder
	validator *validator.Validate
}

func newNotificationContextParser(decoder *schema.Decoder, validator *validator.Validate) *notificationContextParser {
	return &notificationContextParser{
		decoder:   decoder,
		validator: validator,
	}
}

func (p *notificationContextParser) parse(r *http.Request) (model.INotificationContext, error) {
	var params notificationTypeParams

	decodeErr := p.decodeAndValidate(&params, r)
	if decodeErr != nil {
		return nil, decodeErr
	}

	switch params.Type {
	case model.NotificationTypeUsersConnected:
		var params usersConnectedNotificationContext
		decodeErr := p.decodeAndValidate(&params, r)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return params.toContext(), nil
	case model.NotificationTypeUsersLeftChannel:
		var params usersLeftChannelNotificationContext
		decodeErr := p.decodeAndValidate(&params, r)
		if decodeErr != nil {
			return nil, decodeErr
		}
		return params.toContext(), nil
	default:
		return nil, &ErrUnsupportedType{Type: params.Type}
	}
}

func (p *notificationContextParser) decodeAndValidate(
	notificationContext interface{},
	r *http.Request,
) error {
	decodeErr := p.decoder.Decode(notificationContext, r.URL.Query())
	if decodeErr != nil {
		return fmt.Errorf("decode: %w", decodeErr)
	}

	validateErr := p.validator.Struct(notificationContext)
	if validateErr != nil {
		return fmt.Errorf("validation: %w", validateErr)
	}

	return nil
}
