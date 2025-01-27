-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_stats
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    damage INTEGER DEFAULT 1,
    critical_damage INTEGER DEFAULT 0,
    critical_chance DECIMAL(5, 2) DEFAULT 0,
    gold_multiplier DECIMAL(5, 2) DEFAULT 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_stats;
-- +goose StatementEnd
