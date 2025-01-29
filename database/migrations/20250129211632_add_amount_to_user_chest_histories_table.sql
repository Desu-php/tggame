-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_chest_histories
    ADD COLUMN amount INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_chest_histories
DROP
COLUMN IF EXISTS amount;
-- +goose StatementEnd
