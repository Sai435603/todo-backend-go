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

func (r *TodoRepository) GetAll(ctx context.Context) ([]sqlc.Todo, error) {
	return r.queries.GetTodos(ctx)
}

func (r *TodoRepository) GetByID(ctx context.Context, id int64) (sqlc.Todo, error) {
	return r.queries.GetTodoById(ctx, id)
}

func (r *TodoRepository) Update(ctx context.Context, arg sqlc.UpdateTodoParams) (sqlc.Todo, error) {
	return r.queries.UpdateTodo(ctx, arg)
}

func (r *TodoRepository) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteTodo(ctx, id)
}

func (r *TodoRepository) GetCompleted(ctx context.Context) ([]sqlc.Todo, error) {
	return r.queries.GetCompletedTodos(ctx)
}

func (r *TodoRepository) GetPending(ctx context.Context) ([]sqlc.Todo, error) {
	return r.queries.GetPendingTodos(ctx)
}

func (r *TodoRepository) MarkCompleted(ctx context.Context, id int64) (sqlc.Todo, error) {
	return r.queries.MarkTodoAsCompleted(ctx, id)
}

func (r *TodoRepository) MarkPending(ctx context.Context, id int64) (sqlc.Todo, error) {
	return r.queries.MarkTodoAsPending(ctx, id)
}

func (r *TodoRepository) GetByDateRange(ctx context.Context, arg sqlc.GetTodosByDateRangeParams) ([]sqlc.Todo, error) {
	return r.queries.GetTodosByDateRange(ctx, arg)
}
func (r *TodoRepository) Search(ctx context.Context, query string) ([]sqlc.Todo, error) {
	return r.queries.SearchTodos(ctx, pgtype.Text{
		String: query,
		Valid:  true,
	})
}
