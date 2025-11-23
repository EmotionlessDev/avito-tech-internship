package team

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/domain/team"
)

var errNilTx = fmt.Errorf("transaction is nil")

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) GetMembers(ctx context.Context, tx *sql.Tx, teamName string, excludeIDs []string) ([]team.User, error) {
	if tx == nil {
		return nil, errNilTx
	}

	query := `
        SELECT id, name, team_name, is_active
        FROM users
        WHERE team_name = $1 AND id NOT IN (%s)
    `

	args := make([]any, len(excludeIDs)+1)
	args = append(args, teamName)

	in := ""
	for i := 1; i <= len(excludeIDs); i++ {
		in = fmt.Sprintf("%s, $%d", in, i)
		args = append(args, excludeIDs[i-1])
	}

	in = in[2:]

	rows, err := tx.QueryContext(ctx, fmt.Sprintf(query, in), args)
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
