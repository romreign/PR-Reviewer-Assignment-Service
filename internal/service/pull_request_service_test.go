package service

import (
	"testing"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/inmemory"
)

func TestSelectRandomReviewers(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	userRepo := inmemory.NewUserRepository()
	teamRepo := inmemory.NewTeamRepository()
	service := NewPullRequestService(prRepo, teamRepo, userRepo)

	members := []api.TeamMember{
		{UserId: "u1", Username: "Alice", IsActive: true},
		{UserId: "u2", Username: "Bob", IsActive: true},
		{UserId: "u3", Username: "Charlie", IsActive: true},
		{UserId: "u4", Username: "Diana", IsActive: true},
	}

	reviewers := service.SelectRandomReviewers(members, 2)
	if len(reviewers) != 2 {
		t.Errorf("Expected 2 reviewers, got %d", len(reviewers))
	}

	reviewers = service.SelectRandomReviewers(members[:1], 2)
	if len(reviewers) != 1 {
		t.Errorf("Expected 1 reviewer, got %d", len(reviewers))
	}

	reviewers = service.SelectRandomReviewers([]api.TeamMember{}, 2)
	if len(reviewers) != 0 {
		t.Errorf("Expected 0 reviewers, got %d", len(reviewers))
	}
}

func TestMergePR(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	userRepo := inmemory.NewUserRepository()
	teamRepo := inmemory.NewTeamRepository()
	service := NewPullRequestService(prRepo, teamRepo, userRepo)

	pr := &api.PullRequest{
		PullRequestId:     "pr-1",
		PullRequestName:   "Test PR",
		AuthorId:          "u1",
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: []string{"u2", "u3"},
	}
	err := prRepo.CreatePR(*pr)
	if err != nil {
		return
	}

	mergedPR, err := service.MergePR("pr-1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if mergedPR.Status != api.PullRequestStatusMERGED {
		t.Errorf("Expected status MERGED, got %s", mergedPR.Status)
	}

	if mergedPR.MergedAt == nil {
		t.Error("Expected merged_at to be set")
	}

	mergedPR2, err := service.MergePR("pr-1")
	if err != nil {
		t.Fatalf("Expected no error on second merge, got %v", err)
	}

	if mergedPR2.Status != api.PullRequestStatusMERGED {
		t.Errorf("Expected status MERGED on second merge, got %s", mergedPR2.Status)
	}
}

func TestReassignReviewerForbiddenOnMerged(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	userRepo := inmemory.NewUserRepository()
	teamRepo := inmemory.NewTeamRepository()
	service := NewPullRequestService(prRepo, teamRepo, userRepo)

	pr := &api.PullRequest{
		PullRequestId:     "pr-1",
		PullRequestName:   "Test PR",
		AuthorId:          "u1",
		Status:            api.PullRequestStatusMERGED,
		AssignedReviewers: []string{"u2"},
	}
	prRepo.CreatePR(*pr)

	_, _, err := service.ReassignReviewer("pr-1", "u2")
	if err == nil {
		t.Fatal("Expected error when reassigning reviewer on merged PR")
	}
}

func TestReassignReviewerNotAssigned(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	userRepo := inmemory.NewUserRepository()
	teamRepo := inmemory.NewTeamRepository()
	service := NewPullRequestService(prRepo, teamRepo, userRepo)

	pr := &api.PullRequest{
		PullRequestId:     "pr-1",
		PullRequestName:   "Test PR",
		AuthorId:          "u1",
		Status:            api.PullRequestStatusOPEN,
		AssignedReviewers: []string{"u2"},
	}
	err := prRepo.CreatePR(*pr)
	if err != nil {
		return
	}

	_, _, err = service.ReassignReviewer("pr-1", "u3")
	if err == nil {
		t.Fatal("Expected error when reassigning reviewer who is not assigned")
	}
}
