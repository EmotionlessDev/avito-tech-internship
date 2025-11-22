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
	CreateOrUpdate(ctx context.Context, user *models.User) error
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) UpdateIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	query := `
		UPDATE users
		SET is_active = $1
		WHERE user_id = $2
		RETURNING user_id, username, team_name, is_active
	`

	row := r.db.QueryRowContext(ctx, query, isActive, userID)

	var user models.User
	if err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("update is_active: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, userID string) (*models.User, error) {
	query := `
		SELECT user_id, username, team_name, is_active
		FROM users
		WHERE user_id = $1
	`

	row := r.db.QueryRowContext(ctx, query, userID)

	var user models.User
	if err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepo) CreateOrUpdate(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id)
		DO UPDATE SET 
			username = EXCLUDED.username,
			team_name = EXCLUDED.team_name,
			is_active = EXCLUDED.is_active
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.UserID,
		user.Username,
		user.TeamName,
		user.IsActive,
	)
	if err != nil {
		return fmt.Errorf("create or update user: %w", err)
	}

	return nil
}
