-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_tasks
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    task_id      INTEGER NOT NULL REFERENCES tasks (id) ON DELETE CASCADE,
    completed_at TIMESTAMP
);
CREATE INDEX idx_user_tasks_user_id ON user_tasks (user_id);
CREATE INDEX idx_user_tasks_task_id ON user_tasks (task_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_tasks
-- +goose StatementEnd
