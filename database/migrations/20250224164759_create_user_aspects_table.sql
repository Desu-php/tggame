-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_aspects
(
    id              SERIAL PRIMARY KEY,
    user_id         INTEGER       NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    aspect_id       INTEGER       NOT NULL REFERENCES aspects (id) ON DELETE CASCADE,
    aspect_stat_id  INTEGER       NOT NULL REFERENCES aspect_stats (id) ON DELETE CASCADE,
    level           INTEGER       NOT NULL,
    damage          INTEGER       NOT NULL,
    critical_damage INTEGER       NOT NULL,
    critical_chance DECIMAL(5, 2) NOT NULL,
    gold_multiplier DECIMAL(5, 2) NOT NULL,
    amount          INTEGER       NOT NULL,
    created_at      TIMESTAMP,
    updated_at      TIMESTAMP
);
CREATE INDEX idx_user_aspects_aspect_id ON user_aspects (aspect_id);
CREATE INDEX idx_user_aspects_user_id ON user_aspects (user_id);
CREATE INDEX idx_user_aspects_aspect_stat_id ON user_aspects (aspect_stat_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_aspects;
-- +goose StatementEnd
