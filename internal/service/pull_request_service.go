package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/api"
	"github.com/romreign/PR-Reviewer-Assignment-Service/internal/repository"
)

type PullRequestService struct {
	pullRequestRepository repository.PullRequestRepository
	teamRepository        repository.TeamRepository
	userRepository        repository.UserRepository
}

func NewPullRequestService(
	pullRequestRepository repository.PullRequestRepository,
	teamRepository repository.TeamRepository,
	userRepository repository.UserRepository,
) *PullRequestService {
	return &PullRequestService{
		pullRequestRepository: pullRequestRepository,
		teamRepository:        teamRepository,
		userRepository:        userRepository,
	}
}

func randomIndex(n int) (int, error) {
	maxN := big.NewInt(int64(n))
	num, err := rand.Int(rand.Reader, maxN)
	if err != nil {
		return 0, err
	}
	return int(num.Int64()), nil
}

func (s *PullRequestService) GetActiveTeamMembers(authorID string) ([]api.TeamMember, error) {
	author, err := s.userRepository.FindUserByID(authorID)
	if err != nil {
		return nil, fmt.Errorf("author not found")
	}

	if author.TeamName == "" {
		return nil, fmt.Errorf("author has no team")
	}

	members, err := s.teamRepository.FindTeamMembersByName(author.TeamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	var activeMembers []api.TeamMember
	for _, member := range members {
		if member.IsActive && member.UserId != authorID {
			activeMembers = append(activeMembers, member)
		}
	}

	return activeMembers, nil
}

func (s *PullRequestService) SelectRandomReviewers(members []api.TeamMember, count int) []string {
	if count > len(members) {
		count = len(members)
	}
	if count > 2 {
		count = 2
	}

	if count == 0 {
		return []string{}
	}

	indexes := make([]int, len(members))
	for i := 0; i < len(members); i++ {
		indexes[i] = i
	}

	var reviewers []string
	selected := 0
	used := make(map[int]bool)

	for selected < count {
		idx, err := randomIndex(len(members))
		if err != nil {
			break // в реальном коде лучше обработать ошибку
		}
		if !used[idx] {
			reviewers = append(reviewers, members[idx].UserId)
			used[idx] = true
			selected++
		}
	}
	return reviewers
}

func (s *PullRequestService) CreatePR(pr *api.PullRequest) error {
	activeMembers, err := s.GetActiveTeamMembers(pr.AuthorId)
	if err != nil {
		return err
	}

	reviewers := s.SelectRandomReviewers(activeMembers, 2)
	pr.AssignedReviewers = reviewers
	pr.Status = api.PullRequestStatusOPEN
	now := time.Now()
	pr.CreatedAt = &now

	return s.pullRequestRepository.CreatePR(*pr)
}

func (s *PullRequestService) FindPRByID(prID string) (*api.PullRequest, error) {
	return s.pullRequestRepository.FindPRByID(prID)
}

func (s *PullRequestService) MergePR(prID string) (*api.PullRequest, error) {
	pr, err := s.pullRequestRepository.FindPRByID(prID)
	if err != nil {
		return nil, fmt.Errorf("PR not found")
	}

	if pr.Status != api.PullRequestStatusMERGED {
		pr.Status = api.PullRequestStatusMERGED
		now := time.Now()
		pr.MergedAt = &now
		err = s.pullRequestRepository.UpdatePR(*pr)
		if err != nil {
			return nil, err
		}
	}

	return pr, nil
}

func (s *PullRequestService) ReassignReviewer(prID string, oldReviewerID string) (*api.PullRequest, *string, error) {
	pr, err := s.pullRequestRepository.FindPRByID(prID)
	if err != nil {
		return nil, nil, fmt.Errorf("PR not found")
	}

	if pr.Status == api.PullRequestStatusMERGED {
		return nil, nil, fmt.Errorf("cannot reassign on merged PR")
	}

	found := false
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, nil, fmt.Errorf("reviewer is not assigned to this PR")
	}

	oldReviewer, err := s.userRepository.FindUserByID(oldReviewerID)
	if err != nil {
		return nil, nil, fmt.Errorf("reviewer not found")
	}

	members, err := s.teamRepository.FindTeamMembersByName(oldReviewer.TeamName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get team members")
	}

	var candidates []api.TeamMember
	for _, member := range members {
		if !member.IsActive {
			continue
		}
		if member.UserId == oldReviewerID {
			continue
		}
		alreadyReviewer := false
		for _, reviewer := range pr.AssignedReviewers {
			if reviewer == member.UserId {
				alreadyReviewer = true
				break
			}
		}
		if !alreadyReviewer {
			candidates = append(candidates, member)
		}
	}

	if len(candidates) == 0 {
		return nil, nil, fmt.Errorf("no active replacement candidate in team")
	}

	idx, err := randomIndex(len(candidates))
	if err != nil {
		return nil, nil, err
	}
	newReviewer := candidates[idx].UserId

	newReviewers := []string{}
	for _, reviewer := range pr.AssignedReviewers {
		if reviewer != oldReviewerID {
			newReviewers = append(newReviewers, reviewer)
		}
	}
	newReviewers = append(newReviewers, newReviewer)
	pr.AssignedReviewers = newReviewers

	err = s.pullRequestRepository.UpdatePR(*pr)
	if err != nil {
		return nil, nil, err
	}

	return pr, &newReviewer, nil
}

func (s *PullRequestService) FindPRsByReviewer(userID string) ([]api.PullRequest, error) {
	return s.pullRequestRepository.FindPRsByReviewer(userID)
}

func (s *PullRequestService) GetStatistics() (*api.Statistics, error) {
	stats := &api.Statistics{
		TotalAssignments: 0,
		ByUser:           make(map[string]int),
		ByStatus: struct {
			Open   int `json:"open"`
			Merged int `json:"merged"`
		}{},
	}

	prs, err := s.pullRequestRepository.GetAllPRs()
	if err != nil {
		return stats, nil
	}

	for _, pr := range prs {
		for _, reviewer := range pr.AssignedReviewers {
			stats.TotalAssignments++
			stats.ByUser[reviewer]++
		}

		if pr.Status == api.PullRequestStatusOPEN {
			stats.ByStatus.Open++
		} else if pr.Status == api.PullRequestStatusMERGED {
			stats.ByStatus.Merged++
		}
	}

	return stats, nil
}

func (s *PullRequestService) DeactivateUsersAndReassignPRs(teamName string, userIDs []string) (*api.BatchDeactivateResponse, error) {
	team := s.teamRepository.FindTeamByName(teamName)
	if team.TeamName == "" {
		return nil, fmt.Errorf("team not found")
	}

	response := &api.BatchDeactivateResponse{
		DeactivatedCount: 0,
		ReassignedCount:  0,
		Errors: []struct {
			UserID string `json:"user_id"`
			Error  string `json:"error"`
		}{},
	}

	userIDMap := make(map[string]bool)
	for _, userID := range userIDs {
		userIDMap[userID] = true
	}

	var activeReplacements []string
	allUsers, _ := s.userRepository.GetAllUsers()
	for _, user := range allUsers {
		if user.TeamName == teamName && user.IsActive && !userIDMap[user.UserId] {
			activeReplacements = append(activeReplacements, user.UserId)
		}
	}

	if len(activeReplacements) == 0 {
		return nil, fmt.Errorf("no active team members available for reassignment")
	}

	for _, userID := range userIDs {
		err := s.userRepository.UpdateUserStatus(userID, false)
		if err != nil {
			response.Errors = append(response.Errors, struct {
				UserID string `json:"user_id"`
				Error  string `json:"error"`
			}{UserID: userID, Error: "failed to deactivate"})
			continue
		}
		response.DeactivatedCount++

		prs, _ := s.pullRequestRepository.FindPRsByReviewer(userID)
		for _, pr := range prs {
			if pr.Status != api.PullRequestStatusOPEN {
				continue
			}

			var newReviewers []string
			for _, rev := range pr.AssignedReviewers {
				if rev != userID {
					newReviewers = append(newReviewers, rev)
				}
			}

			if len(newReviewers) < 2 && len(activeReplacements) > 0 {
				replacementIndex, err := randomIndex(len(activeReplacements))
				if err != nil {
					continue
				}
				replacement := activeReplacements[replacementIndex]
				newReviewers = append(newReviewers, replacement)
			}

			pr.AssignedReviewers = newReviewers
			err := s.pullRequestRepository.UpdatePR(pr)
			if err != nil {
				return nil, err
			}
			response.ReassignedCount++
		}
	}

	return response, nil
}
