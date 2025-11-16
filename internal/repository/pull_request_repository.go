package repository

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"

type PullRequestRepository interface {
	CreatePR(pr api.PullRequest) error
	FindPRByID(prID string) (*api.PullRequest, error)
	UpdatePR(pr api.PullRequest) error
	FindPRsByReviewer(userID string) ([]api.PullRequest, error)
	GetAllPRs() ([]api.PullRequest, error)
}
