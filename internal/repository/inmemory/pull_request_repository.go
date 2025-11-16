package inmemory

import (
	"fmt"
	"sync"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
)

type PullRequestRepository struct {
	mu  sync.RWMutex
	prs map[string]*api.PullRequest
}

func NewPullRequestRepository() *PullRequestRepository {
	return &PullRequestRepository{
		prs: make(map[string]*api.PullRequest),
	}
}

func (r *PullRequestRepository) CreatePR(pr api.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prs[pr.PullRequestId]; exists {
		return fmt.Errorf("PR already exists")
	}
	r.prs[pr.PullRequestId] = &pr
	return nil
}

func (r *PullRequestRepository) FindPRByID(prID string) (*api.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	pr, ok := r.prs[prID]
	if !ok {
		return nil, fmt.Errorf("PR not found")
	}
	return pr, nil
}

func (r *PullRequestRepository) UpdatePR(pr api.PullRequest) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.prs[pr.PullRequestId]; !ok {
		return fmt.Errorf("PR not found")
	}
	r.prs[pr.PullRequestId] = &pr
	return nil
}

func (r *PullRequestRepository) FindPRsByReviewer(userID string) ([]api.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []api.PullRequest
	for _, pr := range r.prs {
		for _, reviewer := range pr.AssignedReviewers {
			if reviewer == userID {
				result = append(result, *pr)
				break
			}
		}
	}
	return result, nil
}

func (r *PullRequestRepository) GetAllPRs() ([]api.PullRequest, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []api.PullRequest
	for _, pr := range r.prs {
		result = append(result, *pr)
	}
	return result, nil
}
