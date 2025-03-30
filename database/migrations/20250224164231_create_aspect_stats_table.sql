-- +goose Up
-- +goose StatementBegin
CREATE TABLE aspect_stats
(
    id                   SERIAL PRIMARY KEY,
    aspect_id            INTEGER       NOT NULL REFERENCES aspects (id) ON DELETE CASCADE,
    level                INTEGER       NOT NULL,
    damage               INTEGER       NOT NULL,
    critical_damage      INTEGER       NOT NULL,
    critical_chance      DECIMAL(5, 2) NOT NULL,
    gold_multiplier      DECIMAL(5, 2) NOT NULL,
    amount               INTEGER       NOT NULL,
    amount_growth_factor NUMERIC(5, 2) NOT NULL,
    created_at           TIMESTAMP,
    updated_at           TIMESTAMP
);
CREATE INDEX idx_aspect_stats_aspect_id ON aspect_stats (aspect_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS aspect_stats
-- +goose StatementEnd
