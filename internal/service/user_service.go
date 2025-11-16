package service

import (
	"fmt"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUserByID(userID string) (*api.User, error) {
	return s.userRepository.FindUserByID(userID)
}

func (s *UserService) SetUserStatus(userID string, status bool) (*api.User, error) {
	user, err := s.userRepository.FindUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	err = s.userRepository.UpdateUserStatus(userID, status)
	if err != nil {
		return nil, err
	}

	user.IsActive = status
	return user, nil
}
