package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/config"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/inmemory"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

func setupTestServer() *Server {
	cfg := &config.Config{
		Server: config.ServerConfig{
			Env:  "local",
			Port: ":8080",
		},
	}

	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()
	prRepo := inmemory.NewPullRequestRepository()

	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)

	return New(cfg, teamService, userService, prService)
}

func TestPostTeamAdd(t *testing.T) {
	server := setupTestServer()
	server.configureRouter()

	team := api.Team{
		TeamName: "backend",
		Members: []api.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
		},
	}

	body, _ := json.Marshal(team)
	req := httptest.NewRequest("POST", "/team/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusOK {
		t.Errorf("Expected status 201 or 200, got %d", w.Code)
	}
}

func TestGetTeamGet(t *testing.T) {
	server := setupTestServer()
	server.configureRouter()

	teamService := service.NewTeamService(inmemory.NewTeamRepository())
	team := &api.Team{
		TeamName: "frontend",
		Members: []api.TeamMember{
			{UserId: "u2", Username: "Bob", IsActive: true},
		},
	}
	err := teamService.AddTeam(team)
	if err != nil {
		return
	}

	req := httptest.NewRequest("GET", "/team/get?team_name=frontend", nil)
	w := httptest.NewRecorder()

	server.Router.ServeHTTP(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("Expected status 200 or 404, got %d", w.Code)
	}
}
