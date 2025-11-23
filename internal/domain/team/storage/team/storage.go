package team

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/EmotionlessDev/avito-tech-internship/internal/common"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
	"github.com/lib/pq"
)

var errNilTx = fmt.Errorf("transaction is nil")

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

const createSQL = `INSERT INTO team (name) VALUES ($1) RETURNING id`

func (s *Storage) Create(ctx context.Context, tx *sql.Tx, name string) error {
	if tx == nil {
		return errNilTx
	}

	_, err := tx.Exec("INSERT INTO team (name) VALUES ($1)", name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return common.ErrTeamDuplicate
		}
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}

func (s *Storage) CreateMany(ctx context.Context, tx *sql.Tx, users []team.User) error {
	if tx == nil {
		return errNilTx
	}

	values := make([]string, 0, len(users))
	args := make([]any, 0, len(users)*4)

	argPos := 1
	for _, u := range users {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", argPos, argPos+1, argPos+2, argPos+3))
		args = append(args, u.ID, u.Name, u.TeamName, u.IsActive)
		argPos += 4
	}

	query := fmt.Sprintf(`
        INSERT INTO users (id, name, team_name, is_active)
        VALUES %s
        ON CONFLICT (id) DO UPDATE
        SET name = EXCLUDED.name,
            team_name = EXCLUDED.team_name,
            is_active = EXCLUDED.is_active
    `, strings.Join(values, ", "))

	_, err := tx.Exec(query, args...)
	return err
}
