ALTER TABLE settings DROP COLUMN IF EXISTS enable_dingtalk_oauth;
ALTER TABLE settings DROP COLUMN IF EXISTS dingtalk_client_id;
ALTER TABLE settings DROP COLUMN IF EXISTS dingtalk_client_secret;

ALTER TABLE settings ADD COLUMN dingtalk_oauth JSONB DEFAULT '{}'::JSONB;
ALTER TABLE settings ADD COLUMN custom_oauth JSONB DEFAULT '{}'::JSONB;
ALTER TABLE settings ADD COLUMN enable_auto_login BOOLEAN DEFAULT FALSE;