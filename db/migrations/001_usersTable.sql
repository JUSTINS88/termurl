-- +goose Up
CREATE TABLE users (
    machine_key VARCHAR(100) PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;