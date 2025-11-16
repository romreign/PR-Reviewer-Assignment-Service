package repository

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"

type UserRepository interface {
	FindUserByID(userID string) (*api.User, error)
	UpdateUserStatus(userID string, status bool) error
	GetAllUsers() ([]api.User, error)
}
