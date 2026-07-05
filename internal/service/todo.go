package service

import (
	"context"
	"errors"

	"github.com/Sai435603/todo-backend-go/internal/database/sqlc"
	"github.com/Sai435603/todo-backend-go/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
)

type TodoService struct {
	repo *repository.TodoRepository
}

func New(repo *repository.TodoRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// helper to convert int64 userID to pgtype.Int8
func userIDParam(userID int64) pgtype.Int8 {
	return pgtype.Int8{Int64: userID, Valid: true}
}

func (s *TodoService) CreateTodo(ctx context.Context, title string, description string, userID int64) (sqlc.Todo, error) {
	if title == "" {
		return sqlc.Todo{}, errors.New("title is required")
	}
	params := sqlc.CreateTodoParams{
		Title: title,
		Description: pgtype.Text{
			String: description,
			Valid:  description != "",
		},
		Completed: pgtype.Bool{
			Bool:  false,
			Valid: true,
		},
		UserID: userIDParam(userID),
	}
	return s.repo.Create(ctx, params)
}

func (s *TodoService) GetTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.GetAll(ctx, userIDParam(userID))
}

func (s *TodoService) GetTodo(ctx context.Context, id int64) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *TodoService) UpdateTodo(ctx context.Context, id int64, title string, description string, completed bool) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	if title == "" {
		return sqlc.Todo{}, errors.New("title is required")
	}
	params := sqlc.UpdateTodoParams{
		ID:    id,
		Title: title,
		Description: pgtype.Text{
			String: description,
			Valid:  description != "",
		},
		Completed: pgtype.Bool{
			Bool:  completed,
			Valid: true,
		},
	}
	return s.repo.Update(ctx, params)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid todo id")
	}
	return s.repo.Delete(ctx, id)
}

func (s *TodoService) GetCompletedTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.GetCompleted(ctx, userIDParam(userID))
}

func (s *TodoService) GetPendingTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.GetPending(ctx, userIDParam(userID))
}

func (s *TodoService) MarkTodoCompleted(ctx context.Context, id int64) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	return s.repo.MarkCompleted(ctx, id)
}

func (s *TodoService) MarkTodoPending(ctx context.Context, id int64) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	return s.repo.MarkPending(ctx, id)
}

func (s *TodoService) SearchTodos(ctx context.Context, query string, userID int64) ([]sqlc.Todo, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return s.repo.Search(ctx, sqlc.SearchTodosParams{
		UserID: userIDParam(userID),
		Column2: pgtype.Text{
			String: query,
			Valid:  true,
		},
	})
}

func (s *TodoService) GetTodosByDateRange(ctx context.Context, userID int64, from pgtype.Timestamp, to pgtype.Timestamp) ([]sqlc.Todo, error) {
	params := sqlc.GetTodosByDateRangeParams{
		UserID:      userIDParam(userID),
		CreatedAt:   from,
		CreatedAt_2: to,
	}
	return s.repo.GetByDateRange(ctx, params)
}
