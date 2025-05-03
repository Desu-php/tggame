-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_chests
ALTER COLUMN health TYPE BIGINT,
    ALTER COLUMN current_health TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_chests
ALTER COLUMN health TYPE INTEGER,
    ALTER COLUMN current_health TYPE INTEGER;
-- +goose StatementEnd