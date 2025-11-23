package pullrequest

import "time"

type PullRequest struct {
	ID        string
	Name      string
	CreatedAt *time.Time
	Status    string
	MergedAt  *time.Time
	AuthorID  string
}

type PullRequestWithReviewers struct {
	PullRequest
	AssignedReviewers []string
}
