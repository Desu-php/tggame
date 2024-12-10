-- +goose Up
-- +goose StatementBegin
CREATE TABLE rarities (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    drop_weight INTEGER NOT NULL,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rarities;
-- +goose StatementEnd
