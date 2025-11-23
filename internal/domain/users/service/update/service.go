package update

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/users"
)

type UserStorage interface {
	SetActiveByID(ctx context.Context, tx *sql.Tx, id string, isActive bool) (*users.User, error)
}

type Service struct {
	userStorage UserStorage
	db          *sql.DB
}

func NewService(userStorage UserStorage, db *sql.DB) *Service {
	return &Service{
		userStorage: userStorage,
		db:          db,
	}
}

func (s *Service) SetUserActiveByID(ctx context.Context, id string, isActive bool) (*users.User, error) {
	opts := &sql.TxOptions{Isolation: sql.LevelReadCommitted}

	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	user, err := s.userStorage.SetActiveByID(ctx, tx, id, isActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}
