-- +goose Up
-- +goose StatementBegin
CREATE TABLE chests (
    id SERIAL PRIMARY KEY,
    image VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    health INTEGER NOT NULL,
    is_default BOOLEAN NOT NULL,
    growth_factor NUMERIC(5, 2) NOT NULL,
    start_level INTEGER NOT NULL,
    end_level INTEGER NOT NULL,
    rarity_id INTEGER NOT NULL REFERENCES rarities(id) ON DELETE CASCADE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
CREATE INDEX idx_chests_rarity_id ON chests (rarity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chests;
-- +goose StatementEnd
