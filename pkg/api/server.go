package api

import (
	"context"
	"fmt"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pthd-notifications/pkg/domain"
	"syscall"
	"time"
)

type Server struct {
	host  string
	port  int
	debug bool

	service *domain.Service
}

func NewServer(host string, port int, debug bool, service *domain.Service) *Server {
	return &Server{
		port:    port,
		host:    host,
		service: service,
		debug:   debug,
	}
}

func (s *Server) prepareRouter() *gin.Engine {
	if s.debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(sentrygin.New(sentrygin.Options{}))

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	notificationsHndlr := notificationHandler{service: s.service}
	r.GET("/api/v1/notification", notificationsHndlr.Handle)

	return r
}

func (s *Server) Run(ctx context.Context) {
	router := s.prepareRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 10)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("Server exiting")
}
