package chi_api

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"net/http"
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

	decoder   *schema.Decoder
	validator *validator.Validate
}

func NewServer(host string, port int, debug bool, service *domain.Service) *Server {
	return &Server{
		port:    port,
		host:    host,
		service: service,
		debug:   debug,
	}
}

func (s *Server) prepareRouter() http.Handler {
	s.decoder = initializeDecoder()
	s.validator = initializeValidator()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		render.Status(r, http.StatusOK)
	})

	notificationsHandler := newNotificationHandler(s.service, s.decoder, s.validator)

	r.Get("/api/v1/notification", notificationsHandler.Handle)

	return r
}

func (s *Server) Run(ctx context.Context) error {
	router := s.prepareRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().
				Err(err).
				Msgf("ListenAndServe server error")
		}
	}()

	notifiableCtx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-notifiableCtx.Done()

	log.Info().Msg("Stopping server ...")
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	// catching ctx.Done(). timeout of 5 seconds.

	<-ctx.Done()
	log.Info().Msg("Server stopped")

	return nil
}
