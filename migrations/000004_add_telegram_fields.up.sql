ALTER TABLE users
ADD COLUMN telegram_user_id BIGINT UNIQUE,
ADD COLUMN telegram_username VARCHAR(255);

CREATE TABLE IF NOT EXISTS allowed_users (
    telegram_user_id BIGINT PRIMARY KEY
);
