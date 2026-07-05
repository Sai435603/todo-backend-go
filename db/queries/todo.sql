-- name: GetTodos :many
SELECT *
FROM todos WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetTodoById :one
SELECT * FROM todos WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todos (title, description, completed, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateTodo :one
UPDATE todos 
SET title = $1, description = $2, completed = $3
WHERE id = $4
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;

-- name: GetCompletedTodos :many
SELECT * FROM todos WHERE completed = true AND user_id = $1;

-- name: GetPendingTodos :many
SELECT * FROM todos WHERE completed = false AND user_id = $1;

-- name: MarkTodoAsCompleted :one
UPDATE todos 
SET completed = true
WHERE id = $1
RETURNING *;

-- name: MarkTodoAsPending :one
UPDATE todos
SET completed = false
WHERE id = $1
RETURNING *;

-- name: SearchTodos :many
SELECT * FROM todos
WHERE user_id = $1 AND (title ILIKE '%' || $2 || '%' OR description ILIKE '%' || $2 || '%');

-- name: GetTodosByDateRange :many
SELECT * FROM todos
WHERE user_id = $1 AND created_at BETWEEN $2 AND $3;
