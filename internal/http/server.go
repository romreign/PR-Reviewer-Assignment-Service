package http

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func New(config *config.Config, teamService *service.TeamService) *Server {
	return &Server{
		Config:  config,
		Router:  chi.NewRouter(),
		Logger:  setupLogger(config.Env),
		Handler: handler.NewServerHandler(teamService),
	}
}

func (s *Server) Run() error {
	s.configureRouter()
	return http.ListenAndServe(s.Config.Server.Port, s.Router)
}

func (s *Server) configureRouter() {
	s.Router.Use(middleware.DefaultLogger)
	s.Router.Get("/users/getReview", s.Handler.GetTeamGet)
	s.Router.Get("/team/get", s.Handler.GetTeamGet)
	s.Router.Post("/team/add", s.Handler.PostTeamAdd)
	s.Router.Post("/users/setIsActive", s.Handler.PostUsersSetIsActive)
	s.Router.Post("/pullRequest/create", s.Handler.PostPullRequestCreate)
	s.Router.Post("/pullRequest/merge", s.Handler.PostPullRequestMerge)
	s.Router.Post("/pullRequest/reassign", s.Handler.PostPullRequestReassign)
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
