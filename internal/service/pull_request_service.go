package service

import "github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository"

type PullRequestService struct {
	pullRequestRepository repository.PullRequestRepository
}

func NewPullRequestService(pullRequestRepository repository.PullRequestRepository) *PullRequestService {
	return &PullRequestService{
		pullRequestRepository: pullRequestRepository,
	}
}
