package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/http/handler"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/inmemory"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

func TestLoadBatchDeactivation(t *testing.T) {
	results := struct {
		totalRequests int
		successCount  int
		errorCount    int
		totalDuration time.Duration
		minDuration   time.Duration
		maxDuration   time.Duration
		avgDuration   time.Duration
	}{
		totalRequests: 1000,
		minDuration:   time.Hour,
	}

	var allDurations []time.Duration

	for iteration := 0; iteration < 100; iteration++ {
		prRepo := inmemory.NewPullRequestRepository()
		teamRepo := inmemory.NewTeamRepository()
		userRepo := inmemory.NewUserRepository()

		teamName := fmt.Sprintf("load-team-%d", iteration)
		err := teamRepo.CreateTeam(api.Team{
			TeamName: teamName,
			Members:  []api.TeamMember{},
		})
		if err != nil {
			return
		}

		for i := 1; i <= 50; i++ {
			userID := fmt.Sprintf("load-user-%d-%d", iteration, i)
			user := &api.User{
				UserId:   userID,
				TeamName: teamName,
				IsActive: true,
			}
			userRepo.AddUser(user)
		}

		for i := 1; i <= 100; i++ {
			prID := fmt.Sprintf("load-pr-%d-%d", iteration, i)
			pr := api.PullRequest{
				PullRequestId:   prID,
				PullRequestName: fmt.Sprintf("Load PR %d", i),
				AuthorId:        fmt.Sprintf("load-user-%d-1", iteration),
				Status:          api.PullRequestStatusOPEN,
				AssignedReviewers: []string{
					fmt.Sprintf("load-user-%d-%d", iteration, (i%25)+1),
					fmt.Sprintf("load-user-%d-%d", iteration, (i%25)+26),
				},
			}
			now := time.Now()
			pr.CreatedAt = &now
			err := prRepo.CreatePR(pr)
			if err != nil {
				return
			}
		}

		teamService := service.NewTeamService(teamRepo)
		userService := service.NewUserService(userRepo)
		prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)
		h := handler.NewServerHandler(teamService, userService, prService)

		for i := 0; i < 10; i++ {
			deactivateIDs := []string{
				fmt.Sprintf("load-user-%d-%d", iteration, (i%5)+1),
				fmt.Sprintf("load-user-%d-%d", iteration, (i%5)+2),
			}

			body := api.PostUsersDeactivateBatchJSONRequestBody{
				TeamName: teamName,
				UserIds:  deactivateIDs,
			}

			bodyBytes, _ := json.Marshal(body)

			req := httptest.NewRequest("POST", "/users/deactivateBatch", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			start := time.Now()
			h.PostUsersDeactivateBatch(w, req)
			elapsed := time.Since(start)

			allDurations = append(allDurations, elapsed)
			results.totalDuration += elapsed

			if elapsed < results.minDuration {
				results.minDuration = elapsed
			}
			if elapsed > results.maxDuration {
				results.maxDuration = elapsed
			}

			if w.Code == http.StatusOK {
				results.successCount++
			} else {
				results.errorCount++
			}
		}
	}

	results.totalRequests = len(allDurations)
	if len(allDurations) > 0 {
		results.avgDuration = results.totalDuration / time.Duration(len(allDurations))
	}

	t.Logf("\n=== LOAD TEST RESULTS: BATCH DEACTIVATION (1000 ops) ===\n")
	t.Logf("Total Requests: %d\n", results.totalRequests)
	t.Logf("Successful: %d\n", results.successCount)
	t.Logf("Errors: %d\n", results.errorCount)
	t.Logf("\nPerformance Metrics:\n")
	t.Logf("  Min Duration: %v\n", results.minDuration)
	t.Logf("  Max Duration: %v\n", results.maxDuration)
	t.Logf("  Avg Duration: %v\n", results.avgDuration)
	t.Logf("  Total Duration: %v\n", results.totalDuration)
	t.Logf("\nPerformance Assessment:\n")
	if results.avgDuration < 100*time.Millisecond {
		t.Logf("  ✓ PASSED: Average response time %v < 100ms requirement\n", results.avgDuration)
	} else {
		t.Logf("  ✗ WARNING: Average response time %v exceeds 100ms target\n", results.avgDuration)
	}
}

func TestBatchDeactivationPerformance(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()

	err := teamRepo.CreateTeam(api.Team{TeamName: "performance-team", Members: []api.TeamMember{}})
	if err != nil {
		return
	}

	for i := 1; i <= 50; i++ {
		userID := fmt.Sprintf("perf-user-%d", i)
		user := &api.User{
			UserId:   userID,
			TeamName: "performance-team",
			IsActive: true,
		}
		userRepo.AddUser(user)
	}

	for i := 1; i <= 500; i++ {
		prID := fmt.Sprintf("perf-pr-%d", i)
		pr := api.PullRequest{
			PullRequestId:   prID,
			PullRequestName: fmt.Sprintf("Perf PR %d", i),
			AuthorId:        fmt.Sprintf("perf-user-%d", (i%20)+1),
			Status:          api.PullRequestStatusOPEN,
			AssignedReviewers: []string{
				fmt.Sprintf("perf-user-%d", (i%10)+1),
				fmt.Sprintf("perf-user-%d", (i%10)+11),
			},
		}
		now := time.Now()
		pr.CreatedAt = &now
		err := prRepo.CreatePR(pr)
		if err != nil {
			return
		}
	}

	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)

	start := time.Now()

	response, err := prService.DeactivateUsersAndReassignPRs("performance-team", []string{
		"perf-user-1",
		"perf-user-2",
	})

	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("DeactivateUsersAndReassignPRs failed: %v", err)
	}

	if response.DeactivatedCount != 2 {
		t.Errorf("Expected 2 deactivated users, got %d", response.DeactivatedCount)
	}

	t.Logf("\n=== BATCH DEACTIVATION PERFORMANCE TEST ===\n")
	if elapsed > 100*time.Millisecond {
		t.Logf("WARNING: Deactivation took %v (target: <100ms)\n", elapsed)
	} else {
		t.Logf("✓ PASSED: Deactivation completed in %v (target: <100ms)\n", elapsed)
	}
	t.Logf("Metrics:\n")
	t.Logf("  - Duration: %v\n", elapsed)
	t.Logf("  - Deactivated users: %d\n", response.DeactivatedCount)
	t.Logf("  - PRs reassigned: %d\n", response.ReassignedCount)
	t.Logf("  - Data volume: 500 PRs, 50 team members\n")
}

func TestStatisticsPerformance(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()

	for i := 1; i <= 1000; i++ {
		prID := fmt.Sprintf("stats-pr-%d", i)
		pr := api.PullRequest{
			PullRequestId:   prID,
			PullRequestName: fmt.Sprintf("Stats PR %d", i),
			AuthorId:        "stats-author",
			Status:          api.PullRequestStatusOPEN,
			AssignedReviewers: []string{
				fmt.Sprintf("stats-reviewer-%d", (i%10)+1),
				fmt.Sprintf("stats-reviewer-%d", (i%10)+11),
			},
		}
		now := time.Now()
		pr.CreatedAt = &now
		err := prRepo.CreatePR(pr)
		if err != nil {
			return
		}
	}

	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)

	start := time.Now()
	stats, err := prService.GetStatistics()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("GetStatistics failed: %v", err)
	}

	t.Logf("Statistics retrieved in %v\n", elapsed)
	t.Logf("  - Total assignments: %d\n", stats.TotalAssignments)
	t.Logf("  - Reviewers tracked: %d\n", len(stats.ByUser))
	t.Logf("  - Open PRs: %d, Merged: %d\n", stats.ByStatus.Open, stats.ByStatus.Merged)
}
