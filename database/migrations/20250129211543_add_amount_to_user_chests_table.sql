-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_chests
    ADD COLUMN amount INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_chests
DROP
COLUMN IF EXISTS amount;
-- +goose StatementEnd
