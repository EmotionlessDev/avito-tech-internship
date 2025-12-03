package get

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
)

type Service struct {
	teamStorage team.TeamStorage
	userStorage team.UserStorage

	db *sql.DB
}

func NewService(teamStorage team.TeamStorage, userStorage team.UserStorage, db *sql.DB) *Service {
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

	_, err = s.teamStorage.GetByName(ctx, tx, teamName)
	if err != nil {
		if errors.Is(err, common.ErrTeamNotFound) {
			return nil, common.ErrTeamNotFound
		}

		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	members, err := s.userStorage.GetByTeam(ctx, tx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}

	return &TeamWithMembers{
		Team: &team.Team{
			Name: teamName,
		},
		Members: members,
	}, nil
}
