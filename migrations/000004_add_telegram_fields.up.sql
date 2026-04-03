ALTER TABLE users
ADD COLUMN telegram_user_id BIGINT UNIQUE,
ADD COLUMN telegram_username VARCHAR(255);
