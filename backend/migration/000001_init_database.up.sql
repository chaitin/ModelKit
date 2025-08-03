-- 创建 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建模型提供商表
CREATE TABLE IF NOT EXISTS model_providers (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_base VARCHAR(2048) NOT NULL,
    priority INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_model_providers_name ON model_providers (name);

-- 创建模型提供商模型关联表
CREATE TABLE IF NOT EXISTS model_provider_models (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    provider_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_idx_model_provider_models_provider_id_name ON model_provider_models (provider_id, name);

-- 创建模型表
CREATE TABLE IF NOT EXISTS models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID,
    model_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    show_name VARCHAR(255),
    api_base VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_version VARCHAR(255),
    api_header TEXT,
    description TEXT,
    is_internal BOOLEAN DEFAULT false,
    provider VARCHAR(50) NOT NULL,
    context_length INTEGER DEFAULT 4096,
    status VARCHAR(20) DEFAULT 'active',
    capabilities JSONB,
    parameters JSONB,
    pricing JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_model_name_api_base_type ON models (model_name, api_base, model_type);

-- 添加外键约束
ALTER TABLE model_provider_models 
ADD CONSTRAINT fk_model_provider_models_provider_id 
FOREIGN KEY (provider_id) REFERENCES model_providers(id) ON DELETE SET NULL;

-- 插入初始数据
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