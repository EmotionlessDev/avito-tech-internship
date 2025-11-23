package user

import (
	"context"
	"database/sql"
	"fmt"

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

const setActiveByIDSQL = `UPDATE users SET is_active = $1 WHERE id = $2 RETURNING id, name, team_name, is_active;`

func (s *Storage) SetActiveByID(ctx context.Context, tx *sql.Tx, id string, isActive bool) (*users.User, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var (
		u   pgUser
		err error
	)

	err = tx.QueryRow(
		setActiveByIDSQL,
		isActive,
		id,
	).Scan(&u.id, &u.name, &u.teamName, &u.isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to set is_active user by ID: %w", err)
	}

	return pgUserToDomain(u), nil
}

func pgUserToDomain(u pgUser) *users.User {
	return &users.User{
		ID:       u.id,
		Name:     u.name,
		TeamName: u.teamName,
		IsActive: u.isActive,
	}
}
