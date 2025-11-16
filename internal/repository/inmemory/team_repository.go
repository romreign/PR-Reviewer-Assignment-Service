package inmemory

import (
	"fmt"
	"sync"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type TeamRepository struct {
	mu    sync.RWMutex
	teams map[string]api.Team
}

func NewTeamRepository() *TeamRepository {
	return &TeamRepository{
		teams: make(map[string]api.Team),
	}
}

func (r *TeamRepository) CreateTeam(team api.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.teams[team.TeamName]; exists {
		return fmt.Errorf("team already exists")
	}
	r.teams[team.TeamName] = team
	return nil
}

func (r *TeamRepository) UpdateTeam(team api.Team) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.teams[team.TeamName] = team
	return nil
}

func (r *TeamRepository) ExistTeamByName(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.teams[name]
	return exists
}

func (r *TeamRepository) FindTeamByName(name string) api.Team {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.teams[name]
}

func (r *TeamRepository) FindTeamsByUser(userID string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var teamNames []string
	for teamName, team := range r.teams {
		for _, member := range team.Members {
			if member.UserId == userID {
				teamNames = append(teamNames, teamName)
				break
			}
		}
	}
	return teamNames, nil
}

func (r *TeamRepository) FindTeamMembersByName(teamName string) ([]api.TeamMember, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	team, ok := r.teams[teamName]
	if !ok {
		return nil, fmt.Errorf("team not found")
	}
	return team.Members, nil
}
