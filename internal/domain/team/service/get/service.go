package get

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
)

type TeamStorage interface {
	GetByName(ctx context.Context, tx *sql.Tx, name string) (*team.Team, error)
}

type UserStorage interface {
	GetByTeam(ctx context.Context, tx *sql.Tx, teamName string) ([]team.User, error)
}

type Service struct {
	teamStorage TeamStorage
	userStorage UserStorage

	db *sql.DB
}

func NewService(teamStorage TeamStorage, userStorage UserStorage, db *sql.DB) *Service {
	return &Service{
		teamStorage: teamStorage,
		userStorage: userStorage,
		db:          db,
	}
}

type TeamWithMembers struct {
	Team    *team.Team
	Members []team.User
}

func (s *Service) Get(ctx context.Context, teamName string) (*TeamWithMembers, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		tx.Commit()
	}()

	t, err := s.teamStorage.GetByName(ctx, tx, teamName)
	if err != nil {
		if err == common.ErrTeamNotFound {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	members, err := s.userStorage.GetByTeam(ctx, tx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	return &TeamWithMembers{
		Team:    t,
		Members: members,
	}, nil
}
