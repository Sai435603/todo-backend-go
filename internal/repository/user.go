package repository

import (
	"context"

	"github.com/Sai435603/todo-backend-go/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository provides data-access methods for the users table.
type UserRepository struct {
	queries *sqlc.Queries
	db      *pgxpool.Pool
}

// NewUserRepository creates a UserRepository backed by the given connection pool.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

// Upsert inserts a new user or updates the existing record when the google_id
// already exists. This is the primary entry point for OAuth logins.
func (r *UserRepository) Upsert(ctx context.Context, arg sqlc.UpsertUserParams) (sqlc.User, error) {
	return r.queries.UpsertUser(ctx, arg)
}

// GetByID returns a user by their internal ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (sqlc.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

// GetByGoogleID returns a user by their Google account ID.
func (r *UserRepository) GetByGoogleID(ctx context.Context, googleID string) (sqlc.User, error) {
	return r.queries.GetUserByGoogleID(ctx, googleID)
}
