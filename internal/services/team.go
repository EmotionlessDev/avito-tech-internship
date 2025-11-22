package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/EmotionlessDev/avito-tech-internship/internal/models"
	"github.com/EmotionlessDev/avito-tech-internship/internal/repository"
)

type TeamService struct {
	teams repository.TeamRepository
	users repository.UserRepository
}

func NewTeamService(teams repository.TeamRepository, users repository.UserRepository) *TeamService {
	return &TeamService{
		teams: teams,
		users: users,
	}
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (*models.Team, error) {
	return s.teams.GetByName(ctx, name)
}

func (s *TeamService) CreateTeamWithMembers(ctx context.Context, t models.Team) error {
	if err := s.teams.CreateTeam(ctx, t.TeamName); err != nil {
		return err
	}

	for _, u := range t.Members {
		if _, err := s.users.GetByID(ctx, u.UserID); err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				if err = s.users.UpsertUser(ctx, &u); err != nil {
					return fmt.Errorf("failed to create user %s: %w", u.UserID, err)
				}
			} else {
				return err
			}
		}
	}

	if err := s.teams.AddMembers(ctx, t.TeamName, t.Members); err != nil {
		return fmt.Errorf("failed to add members to team: %w", err)
	}

	return nil
}
