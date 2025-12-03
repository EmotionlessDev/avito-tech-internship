package add

import (
	"context"
	"database/sql"
	"fmt"

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

func (s *Service) Add(ctx context.Context, t *team.Team, members []team.User) error {
	opts := &sql.TxOptions{Isolation: sql.LevelRepeatableRead}

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}

		tx.Commit()
	}()

	err = s.teamStorage.Create(ctx, tx, t.Name)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	for i := range members {
		members[i].TeamName = t.Name
	}

	err = s.userStorage.CreateMany(ctx, tx, members)
	if err != nil {
		return fmt.Errorf("failed to create users: %w", err)
	}

	return nil
}
