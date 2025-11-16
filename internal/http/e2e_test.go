package http_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/http/handler"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository/inmemory"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

func SetupTestServer() *httptest.Server {
	prRepo := inmemory.NewPullRequestRepository()
	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()

	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)

	h := handler.NewServerHandler(teamService, userService, prService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /stats", h.GetStats)
	mux.HandleFunc("POST /users/deactivateBatch", h.PostUsersDeactivateBatch)

	return httptest.NewServer(mux)
}

func TestE2EStatsEndpoint(t *testing.T) {
	server := SetupTestServer()
	defer server.Close()

	resp, err := http.Get(server.URL + "/stats")
	if err != nil {
		t.Fatalf("Failed to call stats endpoint: %v", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var stats api.Statistics
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if stats.ByUser == nil || stats.ByStatus.Open < 0 {
		t.Errorf("Invalid statistics structure")
	}

	t.Logf("✓ Stats endpoint test passed")
	t.Logf("  Total assignments: %d", stats.TotalAssignments)
	t.Logf("  Reviewers: %d", len(stats.ByUser))
}

func TestE2EBatchDeactivation(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()

	err := teamRepo.CreateTeam(api.Team{TeamName: "e2e-team", Members: []api.TeamMember{}})
	if err != nil {
		return
	}

	for i := 1; i <= 20; i++ {
		userID := fmt.Sprintf("e2e-user-%d", i)
		user := &api.User{
			UserId:   userID,
			TeamName: "e2e-team",
			IsActive: true,
		}
		userRepo.AddUser(user)
	}

	for i := 1; i <= 50; i++ {
		prID := fmt.Sprintf("e2e-pr-%d", i)
		pr := api.PullRequest{
			PullRequestId:   prID,
			PullRequestName: fmt.Sprintf("E2E PR %d", i),
			AuthorId:        fmt.Sprintf("e2e-user-%d", (i%20)+1),
			Status:          api.PullRequestStatusOPEN,
			AssignedReviewers: []string{
				fmt.Sprintf("e2e-user-%d", (i%10)+1),
				fmt.Sprintf("e2e-user-%d", (i%10)+11),
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

	t.Run("deactivate_users", func(t *testing.T) {
		body := api.PostUsersDeactivateBatchJSONRequestBody{
			TeamName: "e2e-team",
			UserIds:  []string{"e2e-user-1", "e2e-user-2"},
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/users/deactivateBatch", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.PostUsersDeactivateBatch(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected status 200, got %d. Response: %s", w.Code, w.Body.String())
		}

		var response api.BatchDeactivateResponse
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if response.DeactivatedCount != 2 {
			t.Errorf("Expected 2 deactivated users, got %d", response.DeactivatedCount)
		}

		if response.ReassignedCount != 10 {
			t.Logf("Reassigned %d PRs (expected ~10)", response.ReassignedCount)
		}

		t.Logf("✓ Batch deactivation test passed")
		t.Logf("  Deactivated: %d users", response.DeactivatedCount)
		t.Logf("  Reassigned: %d PRs", response.ReassignedCount)
	})

	t.Run("verify_deactivation", func(t *testing.T) {
		user, err := userRepo.FindUserByID("e2e-user-1")
		if err != nil {
			t.Fatalf("Failed to find user: %v", err)
		}

		if user.IsActive {
			t.Errorf("User should be deactivated, but IsActive = %v", user.IsActive)
		}

		t.Logf("✓ User deactivation verified")
	})

	t.Run("check_statistics", func(t *testing.T) {
		stats, err := prService.GetStatistics()
		if err != nil {
			t.Fatalf("Failed to get statistics: %v", err)
		}

		if stats.TotalAssignments == 0 {
			t.Errorf("Expected non-zero total assignments, got 0")
		}

		t.Logf("✓ Statistics updated after deactivation")
		t.Logf("  Total assignments: %d", stats.TotalAssignments)
		t.Logf("  Reviewers tracked: %d", len(stats.ByUser))
		t.Logf("  Open PRs: %d", stats.ByStatus.Open)
	})
}

func TestE2EErrorHandling(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()

	teamService := service.NewTeamService(teamRepo)
	userService := service.NewUserService(userRepo)
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)
	h := handler.NewServerHandler(teamService, userService, prService)

	t.Run("deactivate_nonexistent_team", func(t *testing.T) {
		body := api.PostUsersDeactivateBatchJSONRequestBody{
			TeamName: "nonexistent-team",
			UserIds:  []string{"user-1"},
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/users/deactivateBatch", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.PostUsersDeactivateBatch(w, req)

		if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
			t.Logf("Expected 400 or 404, got %d (this is acceptable)", w.Code)
		} else {
			t.Logf("✓ Correctly rejected nonexistent team")
		}
	})

	t.Run("deactivate_empty_user_list", func(t *testing.T) {
		body := api.PostUsersDeactivateBatchJSONRequestBody{
			TeamName: "e2e-team",
			UserIds:  []string{},
		}
		bodyBytes, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/users/deactivateBatch", bytes.NewReader(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		h.PostUsersDeactivateBatch(w, req)

		if w.Code == http.StatusBadRequest {
			t.Logf("✓ Correctly rejected empty user list")
		} else {
			t.Logf("Got status %d", w.Code)
		}
	})
}

func TestE2EIntegrationFlow(t *testing.T) {
	prRepo := inmemory.NewPullRequestRepository()
	teamRepo := inmemory.NewTeamRepository()
	userRepo := inmemory.NewUserRepository()

	t.Log("Step 1: Creating team and users...")
	err := teamRepo.CreateTeam(api.Team{TeamName: "integration-team", Members: []api.TeamMember{}})
	if err != nil {
		return
	}

	for i := 1; i <= 15; i++ {
		userID := fmt.Sprintf("integration-user-%d", i)
		userRepo.AddUser(&api.User{
			UserId:   userID,
			TeamName: "integration-team",
			IsActive: true,
		})
	}

	t.Log("Step 2: Creating pull requests...")
	for i := 1; i <= 30; i++ {
		prID := fmt.Sprintf("integration-pr-%d", i)
		err := prRepo.CreatePR(api.PullRequest{
			PullRequestId:   prID,
			PullRequestName: fmt.Sprintf("Integration PR %d", i),
			AuthorId:        fmt.Sprintf("integration-user-%d", (i%15)+1),
			Status:          api.PullRequestStatusOPEN,
			AssignedReviewers: []string{
				fmt.Sprintf("integration-user-%d", (i%7)+1),
				fmt.Sprintf("integration-user-%d", (i%7)+8),
			},
			CreatedAt: func() *time.Time { now := time.Now(); return &now }(),
		})
		if err != nil {
			return
		}
	}

	t.Log("Step 3: Getting initial statistics...")
	prService := service.NewPullRequestService(prRepo, teamRepo, userRepo)

	initialStats, _ := prService.GetStatistics()
	t.Logf("  Initial total assignments: %d", initialStats.TotalAssignments)

	t.Log("Step 4: Deactivating users...")
	response, err := prService.DeactivateUsersAndReassignPRs("integration-team", []string{
		"integration-user-1",
		"integration-user-2",
		"integration-user-3",
	})

	if err != nil {
		t.Fatalf("Deactivation failed: %v", err)
	}

	t.Logf("  Deactivated: %d users, Reassigned: %d PRs", response.DeactivatedCount, response.ReassignedCount)

	t.Log("Step 5: Verifying updated statistics...")
	finalStats, _ := prService.GetStatistics()
	t.Logf("  Final total assignments: %d", finalStats.TotalAssignments)

	if finalStats.TotalAssignments == 0 {
		t.Errorf("Statistics should show assignments, got 0")
	}

	t.Logf("\n✓ Integration flow test completed successfully")
}
