package repository

import (
	"context"

	"github.com/Sai435603/todo-backend-go/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepository struct {
	queries *sqlc.Queries
	db      *pgxpool.Pool
}

func New(db *pgxpool.Pool) *TodoRepository {
	return &TodoRepository{
		queries: sqlc.New(db),
		db:      db,
	}
}

func (r *TodoRepository) Create(ctx context.Context, arg sqlc.CreateTodoParams) (sqlc.Todo, error) {
	return r.queries.CreateTodo(ctx, arg)
}

func (r *TodoRepository) GetAll(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return r.queries.GetTodos(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

func (r *TodoRepository) GetByID(ctx context.Context, id int64, userID int64) (sqlc.Todo, error) {
	return r.queries.GetTodoById(ctx, sqlc.GetTodoByIdParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (r *TodoRepository) Update(ctx context.Context, arg sqlc.UpdateTodoParams) (sqlc.Todo, error) {
	return r.queries.UpdateTodo(ctx, arg)
}

func (r *TodoRepository) Delete(ctx context.Context, id int64, userID int64) error {
	return r.queries.DeleteTodo(ctx, sqlc.DeleteTodoParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (r *TodoRepository) GetCompleted(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return r.queries.GetCompletedTodos(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

func (r *TodoRepository) GetPending(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return r.queries.GetPendingTodos(ctx, pgtype.Int8{Int64: userID, Valid: true})
}

func (r *TodoRepository) MarkCompleted(ctx context.Context, id int64, userID int64) (sqlc.Todo, error) {
	return r.queries.MarkTodoAsCompleted(ctx, sqlc.MarkTodoAsCompletedParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (r *TodoRepository) MarkPending(ctx context.Context, id int64, userID int64) (sqlc.Todo, error) {
	return r.queries.MarkTodoAsPending(ctx, sqlc.MarkTodoAsPendingParams{
		ID:     id,
		UserID: pgtype.Int8{Int64: userID, Valid: true},
	})
}

func (r *TodoRepository) GetByDateRange(ctx context.Context, arg sqlc.GetTodosByDateRangeParams) ([]sqlc.Todo, error) {
	return r.queries.GetTodosByDateRange(ctx, arg)
}

func (r *TodoRepository) Search(ctx context.Context, query string, userID int64) ([]sqlc.Todo, error) {
	return r.queries.SearchTodos(ctx, sqlc.SearchTodosParams{
		Column1: pgtype.Text{String: query, Valid: true},
		UserID:  pgtype.Int8{Int64: userID, Valid: true},
	})
}
