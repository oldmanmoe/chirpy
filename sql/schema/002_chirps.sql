-- +goose Up
CREATE TABLE chirps(
    id UUID PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    body TEXT NOT NULL,
    user_id UUID REFERENCES users ON DELETE CASCADE NOT NULL
);

-- +goose Down
DROP TABLE chirps;
