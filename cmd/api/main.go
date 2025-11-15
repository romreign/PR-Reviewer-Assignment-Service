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
	var teamRepository repository.TeamRepository = postgres.NewTeamRepository(postgres.Open(cfg))
	teamService := service.NewTeamService(teamRepository)
	srv := http.New(cfg, teamService)
	srv.Run()
}
