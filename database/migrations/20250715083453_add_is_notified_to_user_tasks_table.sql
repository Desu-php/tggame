-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_tasks ADD COLUMN is_notified BOOLEAN DEFAULT FALSE NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_tasks DROP COLUMN is_notified;
-- +goose StatementEnd