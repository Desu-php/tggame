-- +goose Up
-- +goose StatementBegin
ALTER TABLE tasks
    ADD COLUMN data JSON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
DROP COLUMN data;
-- +goose StatementEnd
