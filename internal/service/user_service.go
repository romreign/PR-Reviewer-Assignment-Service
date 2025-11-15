package service

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository"

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) SetUserStatus(userId string, status bool) error {
	return nil
}
