-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_stat_histories
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    damage INTEGER NOT NULL,
    critical_damage INTEGER NOT NULL,
    critical_chance DECIMAL(5, 2) NOT NULL,
    gold_multiplier DECIMAL(5, 2) NOT NULL,
    passive_damage INTEGER NOT NULL,
    is_upgrade BOOLEAN NOT NULL,
    attributable_type VARCHAR(255) NOT NULL,
    attributable_id INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX idx_user_stat_histories_user_id ON user_stat_histories (user_id);
CREATE INDEX idx_user_stat_histories_attributable_type_attributable_id ON user_stat_histories (attributable_type, attributable_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_stat_histories
-- +goose StatementEnd
