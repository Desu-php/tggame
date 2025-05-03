-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_chest_histories ALTER COLUMN health TYPE BIGINT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_chest_histories ALTER COLUMN health TYPE INTEGER;
-- +goose StatementEnd
