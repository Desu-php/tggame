-- +goose Up
-- +goose StatementBegin
ALTER TABLE chests
    ADD COLUMN max_rarity_id INTEGER DEFAULT NULL REFERENCES rarities(id) ON DELETE CASCADE;
CREATE INDEX idx_chests_max_rarity_id ON chests (max_rarity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_chests_max_rarity_id;
ALTER TABLE chests DROP COLUMN IF EXISTS max_rarity_id;
-- +goose StatementEnd