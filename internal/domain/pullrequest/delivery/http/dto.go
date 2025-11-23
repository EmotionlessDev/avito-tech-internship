package http

type CreatePullRequestRequest struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
	CreatedAt       string `json:"created_at"`
	MergedAt        string `json:"merged_at,omitempty"`
}
