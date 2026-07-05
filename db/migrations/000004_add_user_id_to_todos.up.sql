ALTER TABLE todos ADD COLUMN user_id BIGINT REFERENCES users(id) ON DELETE CASCADE;

CREATE INDEX idx_todos_user_id ON todos (user_id);
