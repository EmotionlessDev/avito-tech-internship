package services

import (
	"context"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
)

type TeamService interface {
	CreateTeamWithMembers(ctx context.Context, t models.Team) error
	GetTeam(ctx context.Context, name string) (*models.Team, error)
}

type teamService struct {
	teams repository.TeamRepository
	users repository.UserRepository
}

func NewTeamService(teams repository.TeamRepository, users repository.UserRepository) TeamService {
	return &teamService{
		teams: teams,
		users: users,
	}
}

func (s *teamService) CreateTeamWithMembers(ctx context.Context, t models.Team) error {
	if err := s.teams.Create(ctx, t.TeamName); err != nil {
		return err
	}

	for _, u := range t.Members {
		u.TeamName = t.TeamName

		if err := s.users.CreateOrUpdate(ctx, &u); err != nil {
			return fmt.Errorf("failed to upsert user %s: %w", u.UserID, err)
		}
	}

	return nil
}

func (s *teamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	return s.teams.GetByName(ctx, name)
}
