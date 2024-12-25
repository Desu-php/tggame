-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_items_type_id ON items (type_id);
CREATE INDEX idx_items_rarity_id ON items (rarity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_items_type_id;

DROP INDEX IF EXISTS idx_items_rarity_id;
-- +goose StatementEnd
