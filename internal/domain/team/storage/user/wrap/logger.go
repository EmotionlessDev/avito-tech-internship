package wrap

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
)

type LogStorage struct {
	base team.UserStorage
	log  *slog.Logger
}

func NewLogStorage(base team.UserStorage, log *slog.Logger) *LogStorage {
	return &LogStorage{
		base: base,
		log:  log,
	}
}

func (s *LogStorage) CreateMany(ctx context.Context, tx *sql.Tx, users []team.User) error {
	const op = "user.wrap.LogStorage.CreateMany"
	err := s.base.CreateMany(ctx, tx, users)
	if err != nil {
		s.log.ErrorContext(ctx, op+" failed", slog.Int("user_count", len(users)), slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (s *LogStorage) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*team.User, error) {
	const op = "user.wrap.LogStorage.GetByID"
	user, err := s.base.GetByID(ctx, tx, id)
	if err != nil {
		s.log.ErrorContext(ctx, op+" failed", slog.Int64("user_id", id), slog.String("error", err.Error()))
		return user, err
	}
	return user, nil
}

func (s *LogStorage) GetByTeam(ctx context.Context, tx *sql.Tx, teamName string) ([]team.User, error) {
	const op = "user.wrap.LogStorage.GetByTeam"
	users, err := s.base.GetByTeam(ctx, tx, teamName)
	if err != nil {
		s.log.ErrorContext(ctx, op+" failed", slog.String("team_name", teamName), slog.String("error", err.Error()))
		return users, err
	}
	return users, nil
}
