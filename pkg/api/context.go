package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"pthd-notifications/pkg/domain/model"
	"strings"
)

type notificationTypeParams struct {
	Type string `form:"type" bidning:"required"`
}

type usersConnectedNotificationContext struct {
	DiscordId int64  `form:"discord_id" binding:"required"`
	Usernames string `form:"usernames" binding:"required"`
	Type      string `form:"type" binding:"required"`
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
	DiscordId int64  `form:"discord_id" binding:"required"`
	Type      string `form:"type" binding:"required"`
}

func (np *usersLeftChannelNotificationContext) toContext() model.INotificationContext {

	return &model.UsersLeftChannelNotificationContext{
		Id: np.DiscordId,
	}
}

func parseNotificationContext(c *gin.Context) (model.INotificationContext, error) {
	var params notificationTypeParams
	bindErr := c.Bind(&params)
	if bindErr != nil {
		return nil, bindErr
	}

	switch params.Type {
	case model.UsersConnectedType:
		var params usersConnectedNotificationContext
		bindErr := c.Bind(&params)
		if bindErr != nil {
			return nil, bindErr
		}
		return params.toContext(), nil
	case model.UsersLeftChannelType:
		var params usersLeftChannelNotificationContext
		bindErr := c.Bind(&params)
		if bindErr != nil {
			return nil, bindErr
		}
		return params.toContext(), nil
	default:
		return nil, errors.New("unsupported notification type")
	}
}
