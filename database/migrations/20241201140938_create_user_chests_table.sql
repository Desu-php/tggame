-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_chests (
    id SERIAL PRIMARY KEY,               -- Уникальный идентификатор
    user_id INTEGER NOT NULL,            -- Ссылка на пользователя
    chest_id INTEGER NOT NULL,           -- Ссылка на сундук
    current_health INTEGER NOT NULL,     -- Текущее здоровье сундука
    level INTEGER NOT NULL,             -- Текущий уровень сундука
    created_at TIMESTAMP DEFAULT NOW(),  -- Дата создания записи
    updated_at TIMESTAMP DEFAULT NOW(),  -- Дата обновления записи
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (chest_id) REFERENCES chests (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_chests;
-- +goose StatementEnd
