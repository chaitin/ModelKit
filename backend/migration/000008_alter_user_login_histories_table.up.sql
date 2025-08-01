ALTER TABLE user_login_histories DROP COLUMN IF EXISTS device;

ALTER TABLE user_login_histories ADD COLUMN IF NOT EXISTS hostname VARCHAR(255);
ALTER TABLE user_login_histories ADD COLUMN IF NOT EXISTS os_type VARCHAR(255);
ALTER TABLE user_login_histories ADD COLUMN IF NOT EXISTS os_release VARCHAR(255);
ALTER TABLE user_login_histories ADD COLUMN IF NOT EXISTS client_id VARCHAR(255);