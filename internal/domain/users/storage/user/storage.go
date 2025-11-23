package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/users"
)

var errNilTx = fmt.Errorf("transaction is nil")

type pgUser struct {
	id       string `db:"id"`
	name     string `db:"name"`
	teamName string `db:"team_name"`
	isActive bool   `db:"is_active"`
}

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

const getByIDSQL = `SELECT id, name, team_name, is_active FROM users WHERE id = $1`

func (s *Storage) GetByID(ctx context.Context, tx *sql.Tx, id int64) (*users.User, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var (
		u   pgUser
		err error
	)

	err = tx.QueryRow(
		getByIDSQL,
		id,
	).Scan(&u)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return pgUserToDomain(u), nil
}

func pgUserToDomain(u pgUser) *team.User {
	return &team.User{
		ID:       u.id,
		Name:     u.name,
		TeamName: u.teamName,
		IsActive: u.isActive,
	}
}
