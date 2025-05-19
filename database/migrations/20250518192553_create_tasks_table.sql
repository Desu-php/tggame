-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    type         VARCHAR(50)  NOT NULL,
    target_value INTEGER      NOT NULL,
    amount       BIGINT       NOT NULL,
    created_at timestamp not null,
    updated_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks
-- +goose StatementEnd
