package user

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
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

const createManySQL = `
	INSERT INTO
		users (id, name, team_name, is_active)
	VALUES
		%s
	ON CONFLICT (id) DO
	UPDATE
	SET
		name = EXCLUDED.name,
		team_name = EXCLUDED.team_name,
		is_active = EXCLUDED.is_active
`

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

const getByIDSQL = `SELECT id, name, team_name, is_active FROM users WHERE id = $1`

func (s *Storage) GetByID(ctx context.Context, tx *sql.Tx, id string) (*team.User, error) {
	if tx == nil {
		return nil, errNilTx
	}

	var (
		u   pgUser
		err error
	)

	err = tx.QueryRow(getByIDSQL, id).Scan(&u.id, &u.name, &u.teamName, &u.isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return pgUserToDomain(u), nil
}

const getByTeamSQL = `
	SELECT id, name, team_name, is_active
	FROM users
	WHERE team_name = $1
	ORDER BY name
`

func (s *Storage) GetByTeam(ctx context.Context, tx *sql.Tx, teamName string) ([]team.User, error) {
	if tx == nil {
		return nil, errNilTx
	}

	rows, err := tx.Query(getByTeamSQL, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var result []team.User

	for rows.Next() {
		var u pgUser
		if err := rows.Scan(&u.id, &u.name, &u.teamName, &u.isActive); err != nil {
			return nil, err
		}
		result = append(result, *pgUserToDomain(u))
	}

	return result, nil
}

func (s *Storage) GetTeamMembers(ctx context.Context, tx *sql.Tx, teamName string, excludeID string) ([]team.User, error) {
	if tx == nil {
		return nil, fmt.Errorf("tx is nil")
	}

	query := `
        SELECT id, name, team_name, is_active
        FROM users
        WHERE team_name = $1 AND id <> $2
    `

	rows, err := tx.QueryContext(ctx, query, teamName, excludeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team members: %w", err)
	}
	defer rows.Close()

	members := []team.User{}
	for rows.Next() {
		var u team.User
		if err := rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		members = append(members, u)
	}

	return members, nil
}

func pgUserToDomain(u pgUser) *team.User {
	return &team.User{
		ID:       u.id,
		Name:     u.name,
		TeamName: u.teamName,
		IsActive: u.isActive,
	}
}
