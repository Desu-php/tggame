-- +goose Up
-- +goose StatementBegin
ALTER TABLE rarities
ADD COLUMN sort INTEGER NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rarities
DROP COLUMN IF EXISTS sort;
-- +goose StatementEnd
