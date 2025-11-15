package repository

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"

type TeamRepository interface {
	CreateTeam(team api.Team) error
	UpdateTeam(team api.Team) error
	ExistTeamByName(name string) bool
	FindTeamByName(name string) api.Team
	FindTeamsByUser(userID string) ([]string, error)
	FindTeamMembersByName(teamName string) ([]api.TeamMember, error)
}
