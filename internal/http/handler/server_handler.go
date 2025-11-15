package handler

import (
	"net/http"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/service"
)

type ServerHandler struct {
	teamService *service.TeamService
}

func NewServerHandler(teamService *service.TeamService) *ServerHandler {
	return &ServerHandler{
		teamService: teamService,
	}
}

// (POST /pullRequest/create)
func (h *ServerHandler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {

}

// (POST /pullRequest/merge)
func (h *ServerHandler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {

}

// (POST /pullRequest/reassign)
func (h *ServerHandler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {

}

// (POST /team/add)
func (h *ServerHandler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {

}

// (GET /team/get)
func (h *ServerHandler) GetTeamGet(w http.ResponseWriter, r *http.Request, params api.GetTeamGetParams) {

}

// (GET /users/getReview)
func (h *ServerHandler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params api.GetUsersGetReviewParams) {

}

// (POST /users/setIsActive)
func (h *ServerHandler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {

}
