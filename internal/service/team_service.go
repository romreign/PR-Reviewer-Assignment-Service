package service

import (
	"fmt"

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

func (s *TeamService) GetTeamByName(teamName string) (*api.Team, error) {
	if !s.teamRepository.ExistTeamByName(teamName) {
		return nil, fmt.Errorf("team not found")
	}
	team := s.teamRepository.FindTeamByName(teamName)
	if team.TeamName == "" {
		return nil, fmt.Errorf("team not found")
	}
	return &team, nil
}

func (s *TeamService) AddTeam(team *api.Team) error {
	if s.teamRepository.ExistTeamByName(team.TeamName) {
		return fmt.Errorf("team already exists")
	}
	return s.teamRepository.CreateTeam(*team)
}
