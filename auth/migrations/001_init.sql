-- +goose Up
CREATE TABLE IF NOT EXISTS auth_users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_auth_users_email ON auth_users (email);

-- +goose Down
DROP INDEX IF EXISTS idx_auth_users_email;
DROP TABLE IF EXISTS auth_users;
