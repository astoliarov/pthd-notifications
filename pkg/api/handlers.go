package api

import (
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"pthd-notifications/pkg/domain"
	"pthd-notifications/pkg/domain/entities"
	"strings"
)

type notificationsParams struct {
	DiscordId int64  `form:"discord_id" binding:"required"`
	Usernames string `form:"usernames" binding:"required"`
}

func (np *notificationsParams) toContext() *entities.NotificationContext {
	names := strings.Split(np.Usernames, ",")

	return &entities.NotificationContext{
		NamesJoined: np.Usernames,
		Names:       names,
	}
}

type notificationHandler struct {
	service *domain.Service
}

func (handler *notificationHandler) Handle(c *gin.Context) {
	var params notificationsParams

	bindErr := c.Bind(&params)
	if bindErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error": "Invalid parameters",
			},
		)
		return
	}

	sendErr := handler.service.SendNotification(params.DiscordId, params.toContext())
	if sendErr != nil {
		log.Println(sendErr)
		sentrygin.GetHubFromContext(c).CaptureException(sendErr)
		c.AbortWithStatusJSON(http.StatusInternalServerError,
			gin.H{
				"error": "Internal server error",
			},
		)
	}

	c.Status(http.StatusOK)
}
