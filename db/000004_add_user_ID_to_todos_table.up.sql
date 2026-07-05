alter table todos add column user_id bigserial references users(id) on delete cascade;

CREATE INDEX idx_todos_user_id ON todos (user_id);