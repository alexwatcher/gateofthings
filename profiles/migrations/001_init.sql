-- +goose Up
CREATE TABLE IF NOT EXISTS user_profiles (
    id UUID PRIMARY KEY,
    name TEXT,
    avatar BYTEA CHECK (LENGTH(avatar) <= 131072), -- 128kB
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS user_profiles;
