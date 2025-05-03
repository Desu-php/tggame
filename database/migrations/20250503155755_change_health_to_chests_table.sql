-- +goose Up
-- +goose StatementBegin
ALTER TABLE chests ALTER COLUMN health TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chests ALTER COLUMN health TYPE INTEGER;
-- +goose StatementEnd