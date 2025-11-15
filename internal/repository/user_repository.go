package repository

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"

type UserRepository interface {
	FindUserByID(userID string) (*api.User, error)
	UpdateUserStatus(userId string, status bool) error
}
