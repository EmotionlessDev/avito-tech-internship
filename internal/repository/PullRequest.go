package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrPRNotFound error = errors.New("pull request not found")

type PullRequestRepository interface {
	Create(ctx context.Context, pr *models.PullRequest) error
	GetByID(ctx context.Context, prID string) (*models.PullRequest, error)
	UpdateStatus(ctx context.Context, prID string, status string) (*models.PullRequest, error)
}

type PullRequestRepo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *PullRequestRepo {
	return &PullRequestRepo{db: db}
}

func (r *PullRequestRepo) Create(ctx context.Context, pr *models.PullRequest) error {
	_, err := r.db.ExecContext(ctx,
		`
		INSERT INTO pull_requests
        (pull_request_id, pull_request_name, author_id, status, created_at)
        VALUES ($1,$2,$3,$4,NOW())
		`)
	return err
}

func (r *PullRequestRepo) GetByID(ctx context.Context, prID string) (*models.PullRequest, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at FROM pull_requests WHERE pull_request_id = $1`,
		prID)

	var pr models.PullRequest
	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request by ID: %w", err)
	}

	return &pr, nil
}

func (r *PullRequestRepo) UpdateStatus(ctx context.Context, prID string, status string) (*models.PullRequest, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE pull_requests
		SET status = $1
		WHERE pull_request_id = $2
		RETURNING pull_request_id, pull_request_name, author_id, status, created_at
	`, status, prID)

	var pr models.PullRequest
	err := row.Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, ErrPRNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update pull request status: %w", err)
	}

	return &pr, nil
}
