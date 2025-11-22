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
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) UpdateIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, team_name, is_active
	`, isActive, userID)

	var user models.User
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update user is_active: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, userID string) (*models.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT user_id, username, team_name, is_active FROM users WHERE user_id = $1`, userID)

	var user models.User
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}
