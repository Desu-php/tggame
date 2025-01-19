-- +goose Up
-- +goose StatementBegin
CREATE TABLE referral_users (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    referred_user_id INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (referred_user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_referral_users_user_id ON referral_users (user_id);
CREATE INDEX idx_referral_users_referred_user_id ON referral_users (referred_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS referral_users;
-- +goose StatementEnd
