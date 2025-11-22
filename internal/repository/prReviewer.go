package repository

import (
	"context"
	"database/sql"
)

type PRReviewerRepository interface {
	AddReviewer(ctx context.Context, prID, reviewerID string) error
}

type PRReviewerRepo struct {
	db *sql.DB
}

func (r *PRReviewerRepo) AddReviewer(ctx context.Context, prID, reviewerID string) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
		VALUES ($1, $2)
	`, prID, reviewerID)
	return err
}
