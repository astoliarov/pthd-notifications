package api

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"net/http"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/model"
)

type notificationHandler struct {
	service *domain.Service
}

func (handler *notificationHandler) Handle(c *gin.Context) {
	notificationContext, bindErr := parseNotificationContext(c)
	if bindErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error": "Invalid parameters",
			},
		)
		return
	}

	sendErr := handler.service.SendNotification(notificationContext)
	if sendErr != nil {
		switch sendErr.(type) {
		case *domain.ErrNoSettings:
			c.AbortWithStatusJSON(http.StatusNotFound,
				gin.H{
					"error": "No settings for such parameters",
				},
			)
		case *model.ErrNoMessage:
			c.AbortWithStatusJSON(http.StatusNotFound,
				gin.H{
					"error": "No message for this notification type",
				},
			)
		default:
			sentrygin.GetHubFromContext(c).CaptureException(sendErr)
			c.AbortWithStatusJSON(http.StatusInternalServerError,
				gin.H{
					"error": "Internal server error",
				},
			)
		}
		return
	}

	c.Status(http.StatusOK)
}
