-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_logs
(
    id      SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    type         VARCHAR(255) NOT NULL,
    created_at timestamp not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_logs
-- +goose StatementEnd
