package team

import (
	"context"
	"database/sql"
)

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

type Team struct {
	Name string
}

type TeamStorage interface {
	Create(ctx context.Context, tx *sql.Tx, name string) error
	GetByName(ctx context.Context, tx *sql.Tx, name string) (*Team, error)
	CreateMany(ctx context.Context, tx *sql.Tx, users []User) error
}

type UserStorage interface {
	GetByID(ctx context.Context, tx *sql.Tx, id int64) (*User, error)
	CreateMany(ctx context.Context, tx *sql.Tx, users []User) error
	GetByTeam(ctx context.Context, tx *sql.Tx, teamName string) ([]User, error)
}
