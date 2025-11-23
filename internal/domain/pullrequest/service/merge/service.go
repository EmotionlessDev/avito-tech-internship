package create

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest"
)

type PullRequestStorage interface {
	Merge(ctx context.Context, tx *sql.Tx, id string) (*pullrequest.PullRequest, error)
	GetReviewersByID(ctx context.Context, tx *sql.Tx, id string) ([]string, error)
}

type Service struct {
	db        *sql.DB
	prStorage PullRequestStorage
}

func NewService(db *sql.DB, prStorage PullRequestStorage) *Service {
	return &Service{
		db:        db,
		prStorage: prStorage,
	}
}

func (s *Service) MergePR(ctx context.Context, pr pullrequest.PullRequest) (*pullrequest.PullRequestWithReviewers, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	mergedPR, err := s.prStorage.Merge(ctx, tx, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to set merge pr: %w", err)
	}

	reviewers, err := s.prStorage.GetReviewersByID(ctx, tx, pr.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviewers: %w", err)
	}

	return &pullrequest.PullRequestWithReviewers{
		PullRequest:       *mergedPR,
		AssignedReviewers: reviewers,
	}, nil
}
