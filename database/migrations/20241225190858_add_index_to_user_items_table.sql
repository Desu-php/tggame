-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_user_items_user_id ON user_items (user_id);
CREATE INDEX idx_user_items_item_id ON user_items (item_id);
CREATE INDEX idx_user_items_user_chest_history_id ON user_items (user_chest_history_id)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_items_user_id;
DROP INDEX IF EXISTS idx_user_items_item_id;
DROP INDEX IF EXISTS idx_user_items_user_chest_history_id;
-- +goose StatementEnd
