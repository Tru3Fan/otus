ALTER TABLE users
DROP COLUMN IF EXISTS telegram_user_id,
     DROP COLUMN IF EXIST telegram_username;