package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) FindUserByID(userID string) (*api.User, error) {
	var user api.User
	err := r.db.Get(&user, `
		SELECT user_id as "user_id", username, team_name as "team_name", is_active as "is_active"
		FROM users WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserStatus(userID string, status bool) error {
	result, err := r.db.Exec("UPDATE users SET is_active = $1 WHERE user_id = $2", status, userID)
	if err != nil {
		return fmt.Errorf("failed to update user status: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *UserRepository) GetAllUsers() ([]api.User, error) {
	var users []api.User
	err := r.db.Select(&users, `
		SELECT user_id as "user_id", username, team_name as "team_name", is_active as "is_active"
		FROM users
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}
