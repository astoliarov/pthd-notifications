package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pthd-notifications/pkg/domain"
)

type Server struct {
	host string
	port int

	service *domain.Service
}

func NewServer(host string, port int, service *domain.Service) *Server {
	return &Server{
		port:    port,
		host:    host,
		service: service,
	}
}

func (s *Server) Run() error {
	r := gin.Default()
	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	notificationsHndlr := notificationHandler{service: s.service}
	r.GET("/api/v1/notification", notificationsHndlr.Handle)

	return r.Run(fmt.Sprintf("%s:%d", s.host, s.port))
}
