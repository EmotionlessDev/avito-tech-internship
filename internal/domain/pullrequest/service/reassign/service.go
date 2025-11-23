package create

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"slices"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest"
)

type PullRequestStorage interface {
	Reassign(ctx context.Context, tx *sql.Tx, id string, oldReviewerID, newReviewerID string) error
	GetByID(ctx context.Context, tx *sql.Tx, id string) (*pullrequest.PullRequest, error)
	GetReviewersByID(ctx context.Context, tx *sql.Tx, id string) ([]string, error)
}

type TeamStorage interface {
	GetTeamMembers(ctx context.Context, tx *sql.Tx, teamName string, excludeIDs []string) ([]string, error)
}

type Service struct {
	db          *sql.DB
	prStorage   PullRequestStorage
	teamStorage TeamStorage
}

func NewService(db *sql.DB, prStorage PullRequestStorage, teamStorage TeamStorage) *Service {
	return &Service{
		db:          db,
		prStorage:   prStorage,
		teamStorage: teamStorage,
	}
}

func (s *Service) ReassignPR(ctx context.Context, id string, oldReviewerID string) (*pullrequest.PullRequestWithReviewers, string, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return nil, "", fmt.Errorf("failed to begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	reviewers, err := s.prStorage.GetReviewersByID(ctx, tx, id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get reviewers: %w", err)
	}

	if len(reviewers) == 0 || !slices.Contains(reviewers, oldReviewerID) {
		return nil, "", common.ErrNotAssignedReviewer
	}

	candidates, err := s.teamStorage.GetTeamMembers(ctx, tx, id, reviewers)
	if err != nil {
		return nil, "", err
	}

	if len(candidates) == 0 {
		return nil, "", common.ErrNoCandidates
	}

	n := rand.Int31n(int32(len(candidates)))
	newReviewerID := candidates[n]

	for i := 0; i < len(reviewers); i++ {
		if reviewers[i] == oldReviewerID {
			reviewers[i] = newReviewerID
			break
		}
	}

	err = s.prStorage.Reassign(ctx, tx, id, oldReviewerID, newReviewerID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to reassign pr: %w", err)
	}

	pr, err := s.prStorage.GetByID(ctx, tx, id)
	if err != nil {
		return nil, "", fmt.Errorf("failed to reassign pr: %w", err)
	}

	return &pullrequest.PullRequestWithReviewers{
		PullRequest:       *pr,
		AssignedReviewers: reviewers,
	}, newReviewerID, nil
}
