-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_items ALTER COLUMN user_chest_history_id DROP NOT NULL;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_items
    ALTER COLUMN user_chest_history_id SET NOT NULL;
-- +goose StatementEnd
