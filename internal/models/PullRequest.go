package models

import "time"

type PullRequest struct {
	PullRequestID     string     `db:"pull_request_id" json:"pull_request_id"`
	PullRequestName   string     `db:"pull_request_name" json:"pull_request_name"`
	AuthorID          string     `db:"author_id" json:"author_id"`
	Status            string     `db:"status" json:"status"` // OPEN | MERGED
	AssignedReviewers []string   `json:"assigned_reviewers" db:"-"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	MergedAt          *time.Time `db:"merged_at" json:"mergedAt"`
}
