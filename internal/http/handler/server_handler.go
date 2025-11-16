package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

type ServerHandler struct {
	teamService *service.TeamService
	userService *service.UserService
	prService   *service.PullRequestService
}

func NewServerHandler(teamService *service.TeamService, userService *service.UserService, prService *service.PullRequestService) *ServerHandler {
	return &ServerHandler{
		teamService: teamService,
		userService: userService,
		prService:   prService,
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, code string, message string) {
	response := api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    api.ErrorResponseErrorCode(code),
			Message: message,
		},
	}
	writeJSON(w, status, response)
}

func (h *ServerHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
	var req api.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	err := h.teamService.AddTeam(&req)
	if err != nil {
		if err.Error() == "team already exists" {
			writeError(w, http.StatusBadRequest, "TEAM_EXISTS", "team_name already exists")
		} else {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			slog.Error("Error adding team", "error", err)
		}
		return
	}

	response := map[string]interface{}{
		"team": req,
	}
	writeJSON(w, http.StatusCreated, response)
}

func (h *ServerHandler) GetTeamGet(w http.ResponseWriter, _ *http.Request, params api.GetTeamGetParams) {
	teamName := params.TeamName
	if teamName == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "team_name parameter is required")
		return
	}

	team, err := h.teamService.GetTeamByName(teamName)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "Team not found")
		return
	}

	writeJSON(w, http.StatusOK, team)
}

func (h *ServerHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"user_id"`
		IsActive bool   `json:"is_active"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	user, err := h.userService.SetUserStatus(req.UserID, req.IsActive)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	response := map[string]interface{}{
		"user": user,
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *ServerHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID   string `json:"pull_request_id"`
		PullRequestName string `json:"pull_request_name"`
		AuthorID        string `json:"author_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if prExists, _ := h.prService.FindPRByID(req.PullRequestID); prExists != nil {
		writeError(w, http.StatusConflict, "PR_EXISTS", "PR id already exists")
		return
	}

	pr := &api.PullRequest{
		PullRequestId:   req.PullRequestID,
		PullRequestName: req.PullRequestName,
		AuthorId:        req.AuthorID,
	}

	err := h.prService.CreatePR(pr)
	if err != nil {
		if err.Error() == "author not found" || err.Error() == "author has no team" {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Author or team not found")
		} else {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			slog.Error("Error creating PR", "error", err)
		}
		return
	}

	response := map[string]interface{}{
		"pr": pr,
	}
	writeJSON(w, http.StatusCreated, response)
}

func (h *ServerHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	pr, err := h.prService.MergePR(req.PullRequestID)
	if err != nil {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "PR not found")
		return
	}

	response := map[string]interface{}{
		"pr": pr,
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *ServerHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PullRequestID string `json:"pull_request_id"`
		OldUserID     string `json:"old_user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	pr, newReviewer, err := h.prService.ReassignReviewer(req.PullRequestID, req.OldUserID)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "PR not found" {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "PR not found")
		} else if errMsg == "cannot reassign on merged PR" {
			writeError(w, http.StatusConflict, "PR_MERGED", "cannot reassign on merged PR")
		} else if errMsg == "reviewer is not assigned to this PR" {
			writeError(w, http.StatusConflict, "NOT_ASSIGNED", "reviewer is not assigned to this PR")
		} else if errMsg == "no active replacement candidate in team" {
			writeError(w, http.StatusConflict, "NO_CANDIDATE", "no active replacement candidate in team")
		} else {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
			slog.Error("Error reassigning reviewer", "error", err)
		}
		return
	}

	response := map[string]interface{}{
		"pr":          pr,
		"replaced_by": *newReviewer,
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *ServerHandler) GetUsersGetReview(w http.ResponseWriter, _ *http.Request, params api.GetUsersGetReviewParams) {
	userID := params.UserId
	if userID == "" {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "user_id parameter is required")
		return
	}

	prs, err := h.prService.FindPRsByReviewer(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		slog.Error("Error getting PRs", "error", err)
		return
	}

	var prShorts []api.PullRequestShort
	for _, pr := range prs {
		prShorts = append(prShorts, api.PullRequestShort{
			PullRequestId:   pr.PullRequestId,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorId,
			Status:          api.PullRequestShortStatus(pr.Status),
		})
	}

	response := map[string]interface{}{
		"user_id":       userID,
		"pull_requests": prShorts,
	}
	writeJSON(w, http.StatusOK, response)
}

func (h *ServerHandler) GetStats(w http.ResponseWriter, _ *http.Request) {
	stats, err := h.prService.GetStatistics()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

func (h *ServerHandler) PostUsersDeactivateBatch(w http.ResponseWriter, r *http.Request) {
	var req api.BatchDeactivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if len(req.UserIds) == 0 {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "user_ids cannot be empty")
		return
	}

	result, err := h.prService.DeactivateUsersAndReassignPRs(req.TeamName, req.UserIds)
	if err != nil {
		if err.Error() == "team not found" {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "Team not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	slog.Info("batch deactivation completed", "team", req.TeamName, "deactivated", result.DeactivatedCount, "reassigned", result.ReassignedCount)
	writeJSON(w, http.StatusOK, result)
}
