-- +goose Up
-- +goose StatementBegin
ALTER TABLE rarities
ADD COLUMN craft_chance DECIMAL(5, 2) DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rarities
DROP COLUMN IF EXISTS craft_chance;
-- +goose StatementEnd
