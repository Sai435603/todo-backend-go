CREATE TABLE users (
    id         BIGSERIAL    PRIMARY KEY,
    google_id  VARCHAR(255) NOT NULL UNIQUE,
    email      VARCHAR(255) NOT NULL UNIQUE,
    name       VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_google_id ON users (google_id);
