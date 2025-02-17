-- +goose Up
-- +goose StatementBegin
ALTER TABLE chests
    ADD COLUMN amount INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN amount_growth_factor NUMERIC(5, 2) NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE chests
DROP
COLUMN IF EXISTS amount,
DROP
COLUMN IF EXISTS amount_growth_factor;
-- +goose StatementEnd
