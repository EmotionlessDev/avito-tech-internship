package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrPRNotFound error = errors.New("pull request not found")

type PullRequestRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	Exists(ctx context.Context, prID string) (bool, error)
}

type PullRequestRepo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) Exists(ctx context.Context, prID string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS(SELECT 1 FROM pull_requests WHERE pull_request_id = $1)
	`, prID).Scan(&exists)

	return exists, err
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *models.PullRequest) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status)
		VALUES ($1, $2, $3, 'OPEN')
	`, pr.PullRequestID, pr.PullRequestName, pr.AuthorID)

	return err
}
