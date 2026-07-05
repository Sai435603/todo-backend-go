-- name: GetTodos :many
SELECT * FROM todos
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTodoById :one
SELECT * FROM todos
WHERE id = $1 AND user_id = $2;

-- name: CreateTodo :one
INSERT INTO todos (title, description, completed, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos
SET title = $1, description = $2, completed = $3
WHERE id = $4 AND user_id = $5
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1 AND user_id = $2;

-- name: GetCompletedTodos :many
SELECT * FROM todos
WHERE completed = true AND user_id = $1;

-- name: GetPendingTodos :many
SELECT * FROM todos
WHERE completed = false AND user_id = $1;

-- name: MarkTodoAsCompleted :one
UPDATE todos
SET completed = true
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: MarkTodoAsPending :one
UPDATE todos
SET completed = false
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: SearchTodos :many
SELECT * FROM todos
WHERE (title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%')
  AND user_id = $2;

-- name: GetTodosByDateRange :many
SELECT * FROM todos
WHERE created_at BETWEEN $1 AND $2
  AND user_id = $3;
