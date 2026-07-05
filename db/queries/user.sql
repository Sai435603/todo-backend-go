-- name: UpsertUser :one
INSERT INTO users (google_id, email, name, avatar_url)
VALUES ($1, $2, $3, $4)
ON CONFLICT (google_id) DO UPDATE
    SET email      = EXCLUDED.email,
        name       = EXCLUDED.name,
        avatar_url = EXCLUDED.avatar_url,
        updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByGoogleID :one
SELECT * FROM users WHERE google_id = $1;
