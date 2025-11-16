package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type TeamRepository struct {
	db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) CreateTeam(team api.Team) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	_, err = tx.Exec("INSERT INTO teams (team_name) VALUES ($1)", team.TeamName)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	for _, member := range team.Members {
		_, err = tx.Exec(`
			INSERT INTO users (user_id, username, team_name, is_active) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id) DO UPDATE SET 
				username = $2, team_name = $3, is_active = $4
		`, member.UserId, member.Username, team.TeamName, member.IsActive)
		if err != nil {
			return fmt.Errorf("failed to create/update user: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *TeamRepository) UpdateTeam(team api.Team) error {
	return r.CreateTeam(team)
}

func (r *TeamRepository) ExistTeamByName(name string) bool {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM teams WHERE team_name = $1)", name).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func (r *TeamRepository) FindTeamByName(name string) api.Team {
	team := api.Team{TeamName: name}

	members, err := r.FindTeamMembersByName(name)
	if err != nil {
		return api.Team{}
	}
	team.Members = members
	return team
}

func (r *TeamRepository) FindTeamsByUser(userID string) ([]string, error) {
	var teamNames []string
	err := r.db.Select(&teamNames, `
		SELECT DISTINCT team_name FROM users WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return teamNames, nil
}

func (r *TeamRepository) FindTeamMembersByName(teamName string) ([]api.TeamMember, error) {
	var members []api.TeamMember
	err := r.db.Select(&members, `
		SELECT u.user_id as "user_id", u.username, u.is_active 
		FROM users u
		WHERE u.team_name = $1
		ORDER BY u.user_id
	`, teamName)
	if err != nil {
		return nil, err
	}
	return members, nil
}
