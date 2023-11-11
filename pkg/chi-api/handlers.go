package chi_api

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"net/http"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/model"
)

type notificationHandler struct {
	service *domain.Service

	parser *notificationContextParser
}

func newNotificationHandler(service *domain.Service, decoder *schema.Decoder, validator *validator.Validate) *notificationHandler {
	return &notificationHandler{
		service: service,
		parser:  newNotificationContextParser(decoder, validator),
	}
}

func (handler *notificationHandler) Handle(w http.ResponseWriter, r *http.Request) {
	notificationContext, validationErr := handler.parser.parse(r)
	if validationErr != nil {
		var validationErrors *validator.ValidationErrors
		var unsupportedError *ErrUnsupportedType

		switch {
		case errors.As(validationErr, &validationErrors):
			errorsMap := make(map[string]interface{})
			for _, fieldErr := range validationErr.(validator.ValidationErrors) {
				errorsMap[fieldErr.Field()] = fieldErr.Tag()
			}
			renderError(w, r, validationErr, "Validation error", errorsMap, http.StatusBadRequest)
		case errors.As(validationErr, &unsupportedError):
			errs := map[string]interface{}{"type": "Unsupported type"}
			renderError(w, r, validationErr, "Validation error", errs, http.StatusBadRequest)
		default:
			renderError(w, r, validationErr, "Internal server error", map[string]interface{}{}, http.StatusInternalServerError)
		}
		return
	}

	sendErr := handler.service.SendNotification(notificationContext)
	if sendErr != nil {
		var noSettingsErr *domain.ErrNoSettings
		var noMessageErr *model.ErrNoMessage

		switch {
		case errors.As(sendErr, &noSettingsErr):
			errs := map[string]interface{}{"parameters": "No settings for such parameters"}
			renderError(w, r, sendErr, "Send error", errs, http.StatusBadRequest)
		case errors.As(sendErr, &noMessageErr):
			errs := map[string]interface{}{"parameters": "No message for such parameters"}
			renderError(w, r, sendErr, "Send error", errs, http.StatusBadRequest)
		default:
			renderError(w, r, sendErr, "Send error", map[string]interface{}{}, http.StatusInternalServerError)
		}
		return
	}

	render.Status(r, http.StatusOK)
}

func renderError(w http.ResponseWriter, r *http.Request, err error, reason string, errors map[string]interface{}, status int) {
	log.Debug().Err(err).Interface("errors", errors).Msg(reason)
	render.JSON(w, r, map[string]interface{}{"reason": reason, "errors": errors})
	render.Status(r, status)
}
