-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL,
    model_type VARCHAR(255) NOT NULL,
    model_id INTEGER NOT NULL,
    type SMALLINT NOT NULL,
    old_balance BIGINT NOT NULL,
    new_balance BIGINT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_transactions_user_id ON transactions (user_id);
CREATE INDEX idx_transactions_model_type_model_id ON transactions (model_type, model_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions
-- +goose StatementEnd
