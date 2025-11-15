package repository

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"

type PullRequestRepository interface {
	CreatePR(pr api.PullRequest) error
	FindPPById(prID string) (*api.PullRequest, error)
	UpdatePR(pr api.PullRequest) error
	FindPRsByReviewer(userID string) ([]api.PullRequest, error)
}
