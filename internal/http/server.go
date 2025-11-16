package http

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/config"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/http/handler"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type Server struct {
	Config  *config.Config
	Router  *chi.Mux
	Logger  *slog.Logger
	Handler *handler.ServerHandler
}

func New(config *config.Config, teamService *service.TeamService, userService *service.UserService, prService *service.PullRequestService) *Server {
	return &Server{
		Config:  config,
		Router:  chi.NewRouter(),
		Logger:  setupLogger(config.Server.Env),
		Handler: handler.NewServerHandler(teamService, userService, prService),
	}
}

func (s *Server) Run() error {
	s.configureRouter()
	srv := &http.Server{
		Addr:         s.Config.Server.Port,
		Handler:      s.Router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	return srv.ListenAndServe()
}

func (s *Server) configureRouter() {
	s.Router.Use(middleware.DefaultLogger)

	wrapper := &api.ServerInterfaceWrapper{
		Handler: s.Handler,
	}

	s.Router.Post("/team/add", wrapper.PostTeamAdd)
	s.Router.Get("/team/get", wrapper.GetTeamGet)
	s.Router.Post("/users/setIsActive", wrapper.PostUsersSetIsActive)
	s.Router.Get("/users/getReview", wrapper.GetUsersGetReview)
	s.Router.Post("/pullRequest/create", wrapper.PostPullRequestCreate)
	s.Router.Post("/pullRequest/merge", wrapper.PostPullRequestMerge)
	s.Router.Post("/pullRequest/reassign", wrapper.PostPullRequestReassign)
	s.Router.Get("/stats", s.Handler.GetStats)
	s.Router.Post("/users/deactivateBatch", s.Handler.PostUsersDeactivateBatch)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
