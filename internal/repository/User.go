package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrUserNotFound error = errors.New("user not found")

type UserRepository interface {
	UpdateIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetByID(ctx context.Context, userID string) (*models.User, error)
	UpsertUser(ctx context.Context, user *models.User) error
	ListTeams(ctx context.Context, userID string) ([]string, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) UpsertUser(ctx context.Context, user *models.User) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO users (user_id, username, is_active)
        VALUES ($1, $2, $3)
        ON CONFLICT (user_id) DO UPDATE
        SET username = EXCLUDED.username,
            is_active = EXCLUDED.is_active
    `, user.UserID, user.Username, user.IsActive)
	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}
	return nil
}

func (r *UserRepo) ListTeams(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT team_name FROM team_members WHERE user_id = $1
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []string
	for rows.Next() {
		var teamName string
		if err := rows.Scan(&teamName); err != nil {
			return nil, err
		}
		teams = append(teams, teamName)
	}

	return teams, nil
}

func (r *UserRepo) GetByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, username, is_active FROM users WHERE user_id = $1
	`, userID).Scan(&user.UserID, &user.Username, &user.IsActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) UpdateIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	result, err := r.db.ExecContext(ctx, `
		UPDATE users SET is_active = $1 WHERE user_id = $2
	`, isActive, userID)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, ErrUserNotFound
	}
	return r.GetByID(ctx, userID)
}
