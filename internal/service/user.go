package service

import (
	"context"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/database/sqlc"
	"github.com/Sai435603/todo-backend-go/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserService contains business logic for user operations.
type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a UserService.
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// FindOrCreateUser upserts a user from Google profile data.
// On first login a new row is created; on subsequent logins the profile is refreshed.
func (s *UserService) FindOrCreateUser(ctx context.Context, gu *auth.GoogleUser) (sqlc.User, error) {
	params := sqlc.UpsertUserParams{
		GoogleID: gu.ID,
		Email:    gu.Email,
		Name:     gu.Name,
		AvatarUrl: pgtype.Text{
			String: gu.AvatarURL,
			Valid:  gu.AvatarURL != "",
		},
	}
	return s.repo.Upsert(ctx, params)
}

// GetUser returns a user by internal ID.
func (s *UserService) GetUser(ctx context.Context, id int64) (sqlc.User, error) {
	return s.repo.GetByID(ctx, id)
}
