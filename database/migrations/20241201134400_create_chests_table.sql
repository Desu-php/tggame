-- +goose Up
-- +goose StatementBegin
CREATE TABLE chests (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    health INTEGER NOT NULL,
    is_default BOOLEAN NOT NULL,
    growth_factor NUMERIC(5, 2) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
