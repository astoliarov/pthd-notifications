package api

import (
	"errors"
	"github.com/getsentry/sentry-go"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"net/http"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/model"
)

type notificationHandler struct {
	service domain.IService

	parser *notificationContextParser
}

func newNotificationHandler(service domain.IService, decoder *schema.Decoder, validator *validator.Validate) *notificationHandler {
	return &notificationHandler{
		service: service,
		parser:  newNotificationContextParser(decoder, validator),
	}
}

func (handler *notificationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	hub := sentry.GetHubFromContext(r.Context())

	notificationContext, validationErr := handler.parser.parse(r)
	if validationErr != nil {
		var validationErrors validator.ValidationErrors
		var unsupportedError *ErrUnsupportedType
		switch {
		case errors.As(validationErr, &validationErrors):
			errorsMap := make(map[string]interface{})
			for _, fieldErr := range validationErrors {
				errorsMap[fieldErr.Field()] = fieldErr.Tag()
			}
			renderError(w, r, validationErr, "validation error", errorsMap, http.StatusBadRequest)
		case errors.As(validationErr, &unsupportedError):
			errs := map[string]interface{}{"type": "unsupported type"}
			renderError(w, r, validationErr, "validation error", errs, http.StatusBadRequest)
		default:
			hub.CaptureException(validationErr)
			renderError(w, r, validationErr, "internal server error", map[string]interface{}{}, http.StatusInternalServerError)
		}
		return
	}

	sendErr := handler.service.SendNotification(notificationContext)
	if sendErr != nil {
		var noSettingsErr *domain.ErrNoSettings
		var noMessageErr *model.ErrNoMessage

		switch {
		case errors.As(sendErr, &noSettingsErr):
			errs := map[string]interface{}{"parameters": "no settings for such parameters"}
			renderError(w, r, sendErr, "send error", errs, http.StatusBadRequest)
		case errors.As(sendErr, &noMessageErr):
			errs := map[string]interface{}{"message": "no message for such parameters"}
			renderError(w, r, sendErr, "send error", errs, http.StatusBadRequest)
		default:
			hub.CaptureException(validationErr)
			renderError(w, r, sendErr, "internal server error", map[string]interface{}{}, http.StatusInternalServerError)
		}
		return
	}

	render.Status(r, http.StatusOK)
}

func renderError(w http.ResponseWriter, r *http.Request, err error, reason string, errors map[string]interface{}, status int) {
	log.Debug().Err(err).Interface("errors", errors).Msg(reason)
	render.Status(r, status)
	render.JSON(w, r, map[string]interface{}{"reason": reason, "errors": errors})
}
