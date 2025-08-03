-- 创建 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建模型表
CREATE TABLE IF NOT EXISTS modelkit_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_modelkit_models_provider_model_name ON modelkit_models (provider, model_name);

-- 创建模型API配置表
CREATE TABLE IF NOT EXISTS modelkit_model_api_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_id UUID NOT NULL UNIQUE,
    api_base VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_version VARCHAR(255),
    api_header TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES modelkit_models(id) ON DELETE CASCADE
);

-- 为新表创建索引
CREATE INDEX IF NOT EXISTS idx_modelkit_model_api_configs_model_id ON modelkit_model_api_configs (model_id);
