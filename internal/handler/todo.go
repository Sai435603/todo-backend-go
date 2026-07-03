package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Sai435603/todo-backend-go/internal/response"
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
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	todo, err := h.service.CreateTodo(r.Context(), req.Title, req.Description)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusCreated, todo)
}

func (h *Handler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetTodos(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	todo, err := h.service.GetTodo(r.Context(), id)
	if err != nil {
		_ = response.Error(w, http.StatusNotFound, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return
	}
	var req UpdateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	todo, err := h.service.UpdateTodo(r.Context(), id, req.Title, req.Description, req.Completed)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		return
	}
	if err := h.service.DeleteTodo(r.Context(), id); err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, map[string]string{"message": "todo deleted successfully"})
}

func (h *Handler) GetCompletedTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetCompletedTodos(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) GetPendingTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.GetPendingTodos(r.Context())
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) MarkTodoCompleted(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	todo, err := h.service.MarkTodoCompleted(r.Context(), id)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) MarkTodoPending(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	todo, err := h.service.MarkTodoPending(r.Context(), id)
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todo)
}

func (h *Handler) SearchTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.SearchTodos(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		_ = response.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, todos)
}

func (h *Handler) GetTodosByDateRange(w http.ResponseWriter, r *http.Request) {
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
	res, err := h.service.GetTodosByDateRange(r.Context(), pgtype.Timestamp{Time: from, Valid: true}, pgtype.Timestamp{Time: to, Valid: true})
	if err != nil {
		_ = response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = response.JSON(w, http.StatusOK, res)
}
