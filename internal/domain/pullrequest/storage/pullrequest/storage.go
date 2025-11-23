package pullrequest

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (s *Storage) GetByID(ctx context.Context, tx *sql.Tx, id string) (*pullrequest.PullRequest, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var pr pgPullRequest
	err := tx.QueryRowContext(ctx,
		"SELECT id, name, created_at, status, merged_at, author_id FROM pull_request WHERE id = $1", id,
	).Scan(&pr.id, &pr.name, &pr.createdAt, &pr.status, &pr.mergedAt, &pr.authorID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get PR by id: %w", err)
	}

	return pgPullRequestToDomain(pr), nil
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
