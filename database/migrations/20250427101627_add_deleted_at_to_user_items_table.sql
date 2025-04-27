-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_items
ADD COLUMN deleted_at TIMESTAMP default null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_items
    DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd
