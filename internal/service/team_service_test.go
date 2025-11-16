package service

import (
	"testing"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/inmemory"
)

func TestAddTeam(t *testing.T) {
	repo := inmemory.NewTeamRepository()
	service := NewTeamService(repo)

	team := &api.Team{
		TeamName: "backend",
		Members: []api.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
			{UserId: "u2", Username: "Bob", IsActive: true},
		},
	}

	err := service.AddTeam(team)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	retrievedTeam, err := service.GetTeamByName("backend")
	if err != nil {
		t.Fatalf("Expected to find team, got error: %v", err)
	}

	if retrievedTeam.TeamName != "backend" {
		t.Errorf("Expected team name 'backend', got %s", retrievedTeam.TeamName)
	}

	if len(retrievedTeam.Members) != 2 {
		t.Errorf("Expected 2 members, got %d", len(retrievedTeam.Members))
	}
}

func TestAddTeamDuplicate(t *testing.T) {
	repo := inmemory.NewTeamRepository()
	service := NewTeamService(repo)

	team := &api.Team{
		TeamName: "backend",
		Members: []api.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
		},
	}

	err := service.AddTeam(team)
	if err != nil {
		return
	}

	err = service.AddTeam(team)
	if err == nil {
		t.Fatal("Expected error when creating duplicate team")
	}
}

func TestGetTeamByName(t *testing.T) {
	repo := inmemory.NewTeamRepository()
	service := NewTeamService(repo)

	_, err := service.GetTeamByName("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent team")
	}

	team := &api.Team{
		TeamName: "frontend",
		Members: []api.TeamMember{
			{UserId: "u3", Username: "Charlie", IsActive: false},
		},
	}
	err = service.AddTeam(team)
	if err != nil {
		return
	}

	retrievedTeam, err := service.GetTeamByName("frontend")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedTeam.TeamName != "frontend" {
		t.Errorf("Expected team name 'frontend', got %s", retrievedTeam.TeamName)
	}
}
