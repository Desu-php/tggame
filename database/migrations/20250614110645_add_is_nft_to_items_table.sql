-- +goose Up
-- +goose StatementBegin
ALTER TABLE items
    ADD COLUMN is_nft BOOLEAN DEFAULT false,
    ADD COLUMN quantity INTEGER DEFAULT 2147483647;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE items
    DROP COLUMN is_nft,
    DROP COLUMN quantity;
-- +goose StatementEnd
