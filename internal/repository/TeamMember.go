package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
)

var ErrTeamMemberNotFound error = errors.New("team member not found")

type TeamMemberRepository interface {
	AddMembers(ctx context.Context, teamName string, members []models.User) error
	ListByTeam(ctx context.Context, teamName string) ([]models.User, error)
}

type TeamMemberRepo struct {
	db *sql.DB
}

func NewTeamMemberRepo(db *sql.DB) *TeamMemberRepo {
	return &TeamMemberRepo{db: db}
}

func (r *TeamMemberRepo) AddMembers(ctx context.Context, teamName string, members []models.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	for _, m := range members {
		_, err = tx.ExecContext(ctx, `
            INSERT INTO team_members (team_name, user_id)
            VALUES ($1, $2)
            ON CONFLICT (team_name, user_id) DO NOTHING
        `, teamName, m.UserID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TeamMemberRepo) ListByTeam(ctx context.Context, teamName string) ([]models.User, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT u.user_id, u.username, u.is_active
        FROM team_members tm
        JOIN users u ON tm.user_id = u.user_id
        WHERE tm.team_name = $1
    `, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.UserID, &u.Username, &u.IsActive); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
