package postgres

import (
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

func (r *TeamRepository) ExistTeamByName(name string) bool {
	var team api.Team
	if err := r.db.Get(&team, "SELECT * FROM team WHERE name = $1", name); err != nil {
		return false
	}
	return true
}

func (r *TeamRepository) FindTeamByName(name string) *api.Team {
	var team api.Team
	if err := r.db.Get(&team, "SELECT * FROM team WHERE name = $1", name); err != nil {
		return nil
	}
	return &team
}

func (r *TeamRepository) CreateTeam(team api.Team) error {
	return nil
}
