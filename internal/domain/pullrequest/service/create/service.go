package create

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/pullrequest"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
)

type PullRequestStorage interface {
	Create(ctx context.Context, tx *sql.Tx, pr pullrequest.PullRequest) error
	GetByID(ctx context.Context, tx *sql.Tx, id string) (*pullrequest.PullRequest, error)
}

type UserStorage interface {
	GetTeamMembers(ctx context.Context, tx *sql.Tx, teamName string, excludeID string) ([]team.User, error)
	GetByID(ctx context.Context, tx *sql.Tx, id string) (*team.User, error)
}

type Service struct {
	db          *sql.DB
	prStorage   PullRequestStorage
	userStorage UserStorage
}

func NewService(db *sql.DB, prStorage PullRequestStorage, userStorage UserStorage) *Service {
	return &Service{
		db:          db,
		prStorage:   prStorage,
		userStorage: userStorage,
	}
}

func (s *Service) CreatePR(ctx context.Context, pr pullrequest.PullRequest) (*pullrequest.PullRequestWithReviewers, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
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

	existing, err := s.prStorage.GetByID(ctx, tx, pr.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, common.ErrPRExists
	}
	user, err := s.userStorage.GetByID(ctx, tx, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	members, err := s.userStorage.GetTeamMembers(ctx, tx, user.TeamName, pr.AuthorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	reviewers := []string{}
	if len(members) > 0 {
		n := 2
		if len(members) < 2 {
			n = len(members)
		}
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		perm := rng.Perm(len(members))
		for i := 0; i < n; i++ {
			reviewers = append(reviewers, members[perm[i]].ID)
		}
	}

	pr.Status = "OPEN"
	err = s.prStorage.Create(ctx, tx, pr)
	if err != nil {
		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return &pullrequest.PullRequestWithReviewers{
		PullRequest:       pr,
		AssignedReviewers: reviewers,
	}, nil
}
