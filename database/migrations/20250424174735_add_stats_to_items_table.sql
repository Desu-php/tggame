-- +goose Up
-- +goose StatementBegin
ALTER TABLE items
    ADD COLUMN damage INTEGER DEFAULT 0,
    ADD COLUMN critical_damage INTEGER DEFAULT 0,
    ADD COLUMN critical_chance DECIMAL(5, 2) DEFAULT 0,
    ADD COLUMN gold_multiplier DECIMAL(5, 2) DEFAULT 0,
    ADD COLUMN passive_damage INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items
    DROP COLUMN IF EXISTS damage,
    DROP COLUMN IF EXISTS critical_damage,
    DROP COLUMN IF EXISTS critical_chance,
    DROP COLUMN IF EXISTS gold_multiplier,
    DROP COLUMN IF EXISTS passive_damage;
-- +goose StatementEnd
