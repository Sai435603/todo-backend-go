DROP INDEX IF EXISTS idx_todos_search;
DROP INDEX IF EXISTS idx_todos_completed_false;
DROP INDEX IF EXISTS idx_todos_completed_true;
DROP INDEX IF EXISTS idx_todos_created_at;

DROP EXTENSION IF EXISTS pg_trgm;