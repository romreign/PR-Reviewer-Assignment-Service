package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepositoryImpl(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *UserRepository) FindUserByID(userID string) (*api.User, error) {
	var user api.User
	if err := r.db.Get(&user, "SELECT * FROM users WHERE id = $1", userID); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserStatus(userId string, status bool) error {
	_, err := r.db.Exec("UPDATE users SET is_active=$1 WHERE id=$2", status, userId)
	return err
}
