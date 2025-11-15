package service

import (
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository"
)

type TeamService struct {
	teamRepository repository.TeamRepository
}

func NewTeamService(teamRepository repository.TeamRepository) *TeamService {
	return &TeamService{
		teamRepository: teamRepository,
	}
}

func (s *TeamService) GetTeamByName(teamName string) api.Team {
	return s.teamRepository.FindTeamByName(teamName)
}

func (s *TeamService) AddTeam(teamName string) *api.Team {
	flag := s.teamRepository.ExistTeamByName(teamName)
	if flag {
		return nil
	}
	return nil ////
}
