package pullrequest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest"
)

var errNilTx = fmt.Errorf("transaction is nil")

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

type pgPullRequest struct {
	id        string
	name      string
	createdAt *time.Time
	status    string
	mergedAt  *time.Time
	authorID  string
}

const createSQL = `
	INSERT INTO pull_request (id, name, created_at, author_id)
	VALUES ($1, $2, $3, $4)
`

func (s *Storage) Create(ctx context.Context, tx *sql.Tx, pr pullrequest.PullRequest) error {
	if tx == nil {
		return errNilTx
	}

	_, err := tx.ExecContext(ctx, createSQL, pr.ID, pr.Name, pr.CreatedAt, pr.AuthorID)
	if err != nil {
		return fmt.Errorf("failed to create pull request: %w", err)
	}

	return nil
}

const getByIDSql = `SELECT id, name, created_at, status, merged_at, author_id FROM pull_request WHERE id = $1`

func (s *Storage) GetByID(ctx context.Context, tx *sql.Tx, id string) (*pullrequest.PullRequest, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var pr pgPullRequest
	err := tx.QueryRowContext(ctx,
		getByIDSql, id,
	).Scan(&pr.id, &pr.name, &pr.createdAt, &pr.status, &pr.mergedAt, &pr.authorID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get PR by id: %w", err)
	}

	return pgPullRequestToDomain(pr), nil
}

const getReviewersByIDSql = `SELECT reviewer_id FROM pr_reviwer WHERE request_id = $1`

func (s *Storage) GetReviewersByID(ctx context.Context, tx *sql.Tx, id string) ([]string, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var reviewers []string
	err := tx.QueryRowContext(ctx,
		getReviewersByIDSql, id,
	).Scan(&reviewers)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get PR reviwers by id: %w", err)
	}

	return reviewers, nil
}

const mergeSQL = `
	UPDATE pull_request
	SET status = 'MERGED' merged_at = NOW()
	WHERE id = $1 AND status = 'OPEN'
	RETURNING id, name, created_at, status, merged_at, author_id
`

func (s *Storage) Merge(ctx context.Context, tx *sql.Tx, id string) (*pullrequest.PullRequest, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var pr pgPullRequest
	err := tx.QueryRowContext(ctx, mergeSQL, id).Scan(&pr.id, &pr.name, &pr.createdAt, &pr.status, &pr.mergedAt, &pr.authorID)
	if err == sql.ErrNoRows {
		return nil, common.ErrPRNotFound
	}

	return pgPullRequestToDomain(pr), nil
}

const reaassignSQL = `
	UPDATE pr_reviwer
	SET
		reviewer_id = $3
	WHERE
		reviewer_id = $2
		AND request_id = $1;
`

func (s *Storage) Reassign(ctx context.Context, tx *sql.Tx, id string, old, new string) error {
	if tx == nil {
		return errNilTx
	}

	row, err := tx.ExecContext(ctx, reaassignSQL, id, old, new)
	if err != nil {
		return err
	}

	n, err := row.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to reassign reviewer: %w", err)
	}

	if n == 0 {
		return errors.New("failed to reassign reviewer")
	}

	return nil
}

func pgPullRequestToDomain(pr pgPullRequest) *pullrequest.PullRequest {
	return &pullrequest.PullRequest{
		ID:        pr.id,
		Name:      pr.name,
		CreatedAt: pr.createdAt,
		Status:    pr.status,
		MergedAt:  pr.mergedAt,
		AuthorID:  pr.authorID,
	}
}
