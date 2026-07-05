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

func (s *TodoService) CreateTodo(ctx context.Context, userID int64, title string, description string) (sqlc.Todo, error) {
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
		UserID: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
	}
	return s.repo.Create(ctx, params)
}

func (s *TodoService) GetTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *TodoService) GetTodo(ctx context.Context, userID int64, id int64) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	return s.repo.GetByID(ctx, id, userID)
}

func (s *TodoService) UpdateTodo(ctx context.Context, userID int64, id int64, title string, description string, completed bool) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	if title == "" {
		return sqlc.Todo{}, errors.New("title is required")
	}
	params := sqlc.UpdateTodoParams{
		ID: id,
		Title: title,
		Description: pgtype.Text{
			String: description,
			Valid:  description != "",
		},
		Completed: pgtype.Bool{
			Bool:  completed,
			Valid: true,
		},
		UserID: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
	}
	return s.repo.Update(ctx, params)
}

func (s *TodoService) DeleteTodo(ctx context.Context, userID int64, id int64) error {
	if id <= 0 {
		return errors.New("invalid todo id")
	}
	return s.repo.Delete(ctx, id, userID)
}

func (s *TodoService) GetCompletedTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.GetCompleted(ctx, userID)
}

func (s *TodoService) GetPendingTodos(ctx context.Context, userID int64) ([]sqlc.Todo, error) {
	return s.repo.GetPending(ctx, userID)
}

func (s *TodoService) MarkTodoCompleted(ctx context.Context, userID int64, id int64) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	return s.repo.MarkCompleted(ctx, id, userID)
}

func (s *TodoService) MarkTodoPending(ctx context.Context, userID int64, id int64) (sqlc.Todo, error) {
	if id <= 0 {
		return sqlc.Todo{}, errors.New("invalid todo id")
	}
	return s.repo.MarkPending(ctx, id, userID)
}

func (s *TodoService) SearchTodos(ctx context.Context, userID int64, query string) ([]sqlc.Todo, error) {
	if query == "" {
		return nil, errors.New("search query cannot be empty")
	}
	return s.repo.Search(ctx, query, userID)
}

func (s *TodoService) GetTodosByDateRange(ctx context.Context, userID int64, from pgtype.Timestamp, to pgtype.Timestamp) ([]sqlc.Todo, error) {
	params := sqlc.GetTodosByDateRangeParams{
		CreatedAt:   from,
		CreatedAt_2: to,
		UserID: pgtype.Int8{
			Int64: userID,
			Valid: true,
		},
	}
	return s.repo.GetByDateRange(ctx, params)
}
