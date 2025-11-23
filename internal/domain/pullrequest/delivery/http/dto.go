package http

import "time"

type CreatePullRequestRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	PullRequestID   string     `json:"pull_request_id"`
	PullRequestName string     `json:"pull_request_name"`
	AuthorID        string     `json:"author_id"`
	Status          string     `json:"status"`
	CreatedAt       *time.Time `json:"created_at"`
	MergedAt        *time.Time `json:"merged_at,omitempty"`
}
