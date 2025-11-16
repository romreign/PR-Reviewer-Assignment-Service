package postgres

import (
	"fmt"
	"time"

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
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	var createdAt interface{} = time.Now()
	if pr.CreatedAt != nil {
		createdAt = pr.CreatedAt
	}

	_, err = tx.Exec(`
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, pr.PullRequestId, pr.PullRequestName, pr.AuthorId, pr.Status, createdAt)
	if err != nil {
		return fmt.Errorf("failed to create PR: %w", err)
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err = tx.Exec(`
			INSERT INTO pr_reviewers (pull_request_id, user_id)
			VALUES ($1, $2)
		`, pr.PullRequestId, reviewerID)
		if err != nil {
			return fmt.Errorf("failed to add reviewer: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *PullRequestRepository) FindPRByID(prID string) (*api.PullRequest, error) {
	var pr api.PullRequest
	var createdAt time.Time
	var mergedAt *time.Time

	err := r.db.QueryRow(`
		SELECT pull_request_id as "pull_request_id", pull_request_name as "pull_request_name", author_id as "author_id", status, created_at, merged_at
		FROM pull_requests WHERE pull_request_id = $1
	`, prID).Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &createdAt, &mergedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find PR: %w", err)
	}

	pr.CreatedAt = &createdAt
	pr.MergedAt = mergedAt

	var reviewers []string
	err = r.db.Select(&reviewers, `
		SELECT user_id FROM pr_reviewers WHERE pull_request_id = $1
	`, prID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}
	pr.AssignedReviewers = reviewers

	return &pr, nil
}

func (r *PullRequestRepository) UpdatePR(pr api.PullRequest) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	var mergedAtValue interface{}
	if pr.MergedAt != nil {
		mergedAtValue = pr.MergedAt
	}

	_, err = tx.Exec(`
		UPDATE pull_requests 
		SET status = $1, merged_at = $2
		WHERE pull_request_id = $3
	`, pr.Status, mergedAtValue, pr.PullRequestId)
	if err != nil {
		return fmt.Errorf("failed to update PR: %w", err)
	}

	_, err = tx.Exec("DELETE FROM pr_reviewers WHERE pull_request_id = $1", pr.PullRequestId)
	if err != nil {
		return fmt.Errorf("failed to delete old reviewers: %w", err)
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err = tx.Exec(`
			INSERT INTO pr_reviewers (pull_request_id, user_id)
			VALUES ($1, $2)
		`, pr.PullRequestId, reviewerID)
		if err != nil {
			return fmt.Errorf("failed to add reviewer: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (r *PullRequestRepository) FindPRsByReviewer(userID string) ([]api.PullRequest, error) {
	var prs []api.PullRequest

	rows, err := r.db.Queryx(`
		SELECT pr.pull_request_id as "pull_request_id", pr.pull_request_name as "pull_request_name", pr.author_id as "author_id", pr.status, pr.created_at, pr.merged_at
		FROM pull_requests pr
		WHERE pr.pull_request_id IN (
			SELECT pull_request_id FROM pr_reviewers WHERE user_id = $1
		)
		ORDER BY pr.created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find PRs: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var pr api.PullRequest
		var createdAt time.Time
		var mergedAt *time.Time

		err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &createdAt, &mergedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}

		pr.CreatedAt = &createdAt
		pr.MergedAt = mergedAt

		var reviewers []string
		err = r.db.Select(&reviewers, `
			SELECT user_id FROM pr_reviewers WHERE pull_request_id = $1
		`, pr.PullRequestId)
		if err != nil {
			return nil, fmt.Errorf("failed to get reviewers: %w", err)
		}
		pr.AssignedReviewers = reviewers

		prs = append(prs, pr)
	}

	return prs, rows.Err()
}

func (r *PullRequestRepository) GetAllPRs() ([]api.PullRequest, error) {
	var prs []api.PullRequest

	rows, err := r.db.Queryx(`
		SELECT pull_request_id as "pull_request_id", pull_request_name as "pull_request_name", author_id as "author_id", status, created_at, merged_at
		FROM pull_requests
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to find PRs: %w", err)
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var pr api.PullRequest
		var createdAt time.Time
		var mergedAt *time.Time

		err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &createdAt, &mergedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan PR: %w", err)
		}

		pr.CreatedAt = &createdAt
		pr.MergedAt = mergedAt

		var reviewers []string
		err = r.db.Select(&reviewers, `
			SELECT user_id FROM pr_reviewers WHERE pull_request_id = $1
		`, pr.PullRequestId)
		if err != nil {
			return nil, fmt.Errorf("failed to get reviewers: %w", err)
		}
		pr.AssignedReviewers = reviewers

		prs = append(prs, pr)
	}

	return prs, rows.Err()
}
