package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type PullRequestRepository struct {
	db *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) *PullRequestRepository {
	return &PullRequestRepository{
		db: db,
	}
}

func (r *PullRequestRepository) CreatePR(pr api.PullRequest) error {
	return nil
}

func (r *PullRequestRepository) FindPPById(prID string) (*api.PullRequest, error) {
	var pullRequest api.PullRequest
	if err := r.db.Get(&pullRequest, "SELECT * FROM pull_requests WHERE id = $1", prID); err != nil {
		return nil, err
	}
	return &pullRequest, nil
}

func (r *PullRequestRepository) UpdatePR(pr api.PullRequest) error {
	return nil
}

func (r *PullRequestRepository) FindPRsByReviewer(userID string) ([]api.PullRequest, error) {
	return nil, nil
}
