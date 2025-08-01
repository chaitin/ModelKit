CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    last_active_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admins_username ON admins (username);

CREATE TABLE IF NOT EXISTS invite_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL,
    code VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255),
    password VARCHAR(255),
    email VARCHAR(255),
    avatar_url VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active',
    platform VARCHAR(12),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_users_username ON users (username) WHERE deleted_at IS NULL AND username != '';
CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_users_email ON users (email) WHERE deleted_at IS NULL AND email != '';

CREATE TABLE IF NOT EXISTS user_identities (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    user_id UUID NOT NULL,
    platform VARCHAR(12) NOT NULL,
    identity_id VARCHAR(64) NOT NULL,
    union_id VARCHAR(64),
    nickname VARCHAR(255),
    email VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_user_identities_platform_identity_id ON user_identities (platform, identity_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS user_login_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    ip VARCHAR(255),
    country VARCHAR(255),
    province VARCHAR(255),
    city VARCHAR(255),
    isp VARCHAR(255),
    asn VARCHAR(255),
    client_version VARCHAR(255),
    device VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_user_login_histories_user_id ON user_login_histories (user_id);

CREATE TABLE IF NOT EXISTS admin_login_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL,
    ip VARCHAR(255),
    country VARCHAR(255),
    province VARCHAR(255),
    city VARCHAR(255),
    isp VARCHAR(255),
    asn VARCHAR(255),
    client_version VARCHAR(255),
    device VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_admin_login_histories_admin_id ON admin_login_histories (admin_id);

-- 模型表
CREATE TABLE IF NOT EXISTS models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    api_base VARCHAR(255),
	api_key VARCHAR(255),
	api_version VARCHAR(50),
	description TEXT,
	provider VARCHAR(50) NOT NULL,
	context_length INTEGER DEFAULT 4096,
	status VARCHAR(20) DEFAULT 'active',
	capabilities JSONB,
	parameters JSONB,
	pricing JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_model_api_type ON models (model_name, api_base, model_type);

CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    key VARCHAR(64) NOT NULL UNIQUE,
    name VARCHAR(64) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    last_used TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);

CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    task_id VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    model_id UUID NOT NULL,
    request_id VARCHAR(255),
    model_type VARCHAR(255) NOT NULL,
    completion TEXT,
    is_accept BOOLEAN DEFAULT FALSE,
    program_language VARCHAR(255),
    work_mode VARCHAR(255),
    code_lines BIGINT,
    input_tokens BIGINT,
    output_tokens BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_tasks_task_id ON tasks(task_id);
CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks (user_id);
CREATE INDEX IF NOT EXISTS idx_tasks_model_id ON tasks (model_id);
CREATE INDEX IF NOT EXISTS idx_tasks_model_type ON tasks (model_type);
CREATE INDEX IF NOT EXISTS idx_tasks_work_mode ON tasks (work_mode);
CREATE INDEX IF NOT EXISTS idx_tasks_program_language ON tasks (program_language);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks (created_at);
CREATE INDEX IF NOT EXISTS idx_tasks_updated_at ON tasks (updated_at);

CREATE TABLE IF NOT EXISTS task_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    task_id UUID NOT NULL,
    role VARCHAR(255),
    prompt TEXT,
    completion TEXT,
    output_tokens BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_task_records_task_id ON task_records (task_id);

CREATE TABLE IF NOT EXISTS settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    enable_sso BOOLEAN DEFAULT FALSE,
    force_two_factor_auth BOOLEAN DEFAULT FALSE,
    disable_password_login BOOLEAN DEFAULT FALSE,
    enable_dingtalk_oauth BOOLEAN DEFAULT FALSE,
    dingtalk_client_id VARCHAR(255),
    dingtalk_client_secret VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
