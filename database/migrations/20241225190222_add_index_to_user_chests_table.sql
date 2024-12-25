-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_user_chests_chest_id ON user_chests (chest_id);
CREATE INDEX idx_user_chests_user_id ON user_chests (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_chests_chest_id;

DROP INDEX IF EXISTS idx_user_chests_user_id;
-- +goose StatementEnd
