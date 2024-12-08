-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_chest_histories (
    id SERIAL PRIMARY KEY,               -- Уникальный идентификатор
    user_chest_id INTEGER NOT NULL,           -- Ссылка на сундук
    health INTEGER NOT NULL,             -- Здоровье сундука
    level INTEGER NOT NULL,             -- Текущий уровень сундука
    created_at TIMESTAMP DEFAULT NOW(),  -- Дата создания записи
    updated_at TIMESTAMP DEFAULT NOW(),  -- Дата обновления записи
    FOREIGN KEY (user_chest_id) REFERENCES user_chests (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_chest_histories;
-- +goose StatementEnd