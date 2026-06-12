-- +goose Up
CREATE TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    email TEXT UNIQUE NOT NULL,
    hashed_password TEXT NOT NULL DEFAULT 'unset',
    is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE users CASCADE;