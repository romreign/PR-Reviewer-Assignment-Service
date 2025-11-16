package main

import (
	"log"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/config"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/http"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/postgres"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

func main() {
	cfg, err := config.Load("./config")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := postgres.Open(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer postgres.Close(db)

	var teamRepository repository.TeamRepository = postgres.NewTeamRepository(db)
	var userRepository repository.UserRepository = postgres.NewUserRepository(db)
	var prRepository repository.PullRequestRepository = postgres.NewPullRequestRepository(db)

	teamService := service.NewTeamService(teamRepository)
	userService := service.NewUserService(userRepository)
	prService := service.NewPullRequestService(prRepository, teamRepository, userRepository)

	srv := http.New(cfg, teamService, userService, prService)
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error server run: %v", err)
	}
}
