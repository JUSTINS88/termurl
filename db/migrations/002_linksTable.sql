-- +goose Up
CREATE TABLE links (
    id UUID PRIMARY KEY,
    original_link TEXT NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    clicks INT NOT NULL DEFAULT 0,
    machine_key VARCHAR(100) NOT NULL REFERENCES users(machine_key) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE links;