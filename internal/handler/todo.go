package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Sai435603/todo-backend-go/internal/auth"
	"github.com/Sai435603/todo-backend-go/internal/response"
	"github.com/Sai435603/todo-backend-go/internal/validator"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := validator.ValidateCreateTodo(req.Title, req.Description); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	todo, err := h.service.CreateTodo(r.Context(), userID, req.Title, req.Description)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusCreated, todo)
}

func (h *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	todos, err := h.service.GetTodos(r.Context(), userID)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if errs := validator.ValidateID(id); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	todo, err := h.service.GetTodo(r.Context(), userID, id)
	if err != nil {
		_ = response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := validator.ValidateUpdateTodo(id, req.Title, req.Description); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	todo, err := h.service.UpdateTodo(r.Context(), userID, id, req.Title, req.Description, req.Completed)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if errs := validator.ValidateID(id); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	if err := h.service.DeleteTodo(r.Context(), userID, id); err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, map[string]string{"message": "todo deleted successfully"})
}

func (h *Handler) GetCompletedTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	todos, err := h.service.GetCompletedTodos(r.Context(), userID)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) GetPendingTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	todos, err := h.service.GetPendingTodos(r.Context(), userID)
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) MarkTodoCompleted(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if errs := validator.ValidateID(id); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	todo, err := h.service.MarkTodoCompleted(r.Context(), userID, id)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) MarkTodoPending(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	if errs := validator.ValidateID(id); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	todo, err := h.service.MarkTodoPending(r.Context(), userID, id)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) SearchTodos(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	query := r.URL.Query().Get("q")

	if errs := validator.ValidateSearchQuery(query); errs != nil {
		_ = response.ValidationError(w, errs)
		return
	}

	todos, err := h.service.SearchTodos(r.Context(), userID, query)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) GetTodosByDateRange(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.GetUserID(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusUnauthorized, "not authenticated")
		return
	}

	from, err := time.Parse(time.RFC3339, r.URL.Query().Get("from"))
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid from")
		return
	}
	to, err := time.Parse(time.RFC3339, r.URL.Query().Get("to"))
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid to")
		return
	}
	res, err := h.service.GetTodosByDateRange(r.Context(), userID, pgtype.Timestamp{Time: from, Valid: true}, pgtype.Timestamp{Time: to, Valid: true})
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, res)
}
