package service

import (
	"testing"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/inmemory"
)

func TestGetUserByID(t *testing.T) {
	repo := inmemory.NewUserRepository()
	repo.AddUser(&api.User{UserId: "u1", Username: "Alice", IsActive: true, TeamName: "backend"})
	service := NewUserService(repo)

	user, err := service.GetUserByID("u1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.UserId != "u1" {
		t.Errorf("Expected user ID 'u1', got %s", user.UserId)
	}

	if user.Username != "Alice" {
		t.Errorf("Expected username 'Alice', got %s", user.Username)
	}

	_, err = service.GetUserByID("u999")
	if err == nil {
		t.Fatal("Expected error for non-existent user")
	}
}

func TestSetUserStatus(t *testing.T) {
	repo := inmemory.NewUserRepository()
	repo.AddUser(&api.User{UserId: "u2", Username: "Bob", IsActive: true, TeamName: "backend"})
	service := NewUserService(repo)

	user, err := service.SetUserStatus("u2", false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.IsActive {
		t.Error("Expected user to be inactive")
	}

	user, err = service.SetUserStatus("u2", true)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}
}
