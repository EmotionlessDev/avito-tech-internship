package wrap

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
)

type LogStorage struct {
	base team.TeamStorage
	log  *slog.Logger
}

func NewLogStorage(base team.TeamStorage, log *slog.Logger) *LogStorage {
	return &LogStorage{
		base: base,
		log:  log,
	}
}

func (s *LogStorage) Create(ctx context.Context, tx *sql.Tx, name string) error {
	const op = "wrap.LogStorage.Create"
	err := s.base.Create(ctx, tx, name)
	if err != nil {
		s.log.ErrorContext(ctx, op+" failed", slog.String("team_name", name), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *LogStorage) CreateMany(ctx context.Context, tx *sql.Tx, users []team.User) error {
	const op = "wrap.LogStorage.CreateMany"
	err := s.base.CreateMany(ctx, tx, users)
	if err != nil {
		s.log.ErrorContext(ctx, op+" failed", slog.Int("user_count", len(users)), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *LogStorage) GetByName(ctx context.Context, tx *sql.Tx, name string) (*team.Team, error) {
	const op = "wrap.LogStorage.GetByName"
	teamObj, err := s.base.GetByName(ctx, tx, name)
	if err != nil {
		s.log.ErrorContext(ctx, op+" failed", slog.String("team_name", name), slog.String("error", err.Error()))
		return nil, err
	}
	return teamObj, nil
}
