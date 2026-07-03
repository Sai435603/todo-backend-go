CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE INDEX idx_todos_created_at ON todos (created_at DESC);
CREATE INDEX idx_todos_completed_true ON todos (completed) WHERE completed = true;
CREATE INDEX idx_todos_completed_false ON todos (completed) WHERE completed = false;
CREATE INDEX idx_todos_search ON todos USING GIN (title gin_trgm_ops, description gin_trgm_ops);