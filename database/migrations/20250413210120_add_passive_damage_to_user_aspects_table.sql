-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_aspects
    ADD COLUMN passive_damage INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_aspects
DROP
COLUMN IF EXISTS passive_damage;
-- +goose StatementEnd
