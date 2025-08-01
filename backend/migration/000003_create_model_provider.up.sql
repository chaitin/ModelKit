CREATE TABLE IF NOT EXISTS model_providers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_base VARCHAR(2048) NOT NULL,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_model_providers_name ON model_providers (name);

CREATE TABLE IF NOT EXISTS model_provider_models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    provider_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_model_provider_models_provider_id_name ON model_provider_models (provider_id, name);

INSERT INTO model_providers (id, name, api_base, priority) VALUES
('baizhiyun', '百智云', 'https://model-square.app.baizhi.cloud/v1', 100),
('deepseek', 'DeepSeek', 'https://api.deepseek.com', 90);

INSERT INTO model_provider_models (provider_id, name) VALUES
('baizhiyun', 'deepseek-v3'),
('baizhiyun', 'deepseek-r1'),
('baizhiyun', 'qwen2.5-coder-1.5b-instruct'),
('baizhiyun', 'qwen2.5-coder-3b-instruct'),
('baizhiyun', 'qwen2.5-coder-7b-instruct'),
('deepseek', 'deepseek-chat'),
('deepseek', 'deepseek-reasoner');