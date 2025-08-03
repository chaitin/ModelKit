-- 创建 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建模型表
CREATE TABLE IF NOT EXISTS modelkit_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    provider VARCHAR(50) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    api_base VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_version VARCHAR(255),
    api_header TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_modelkit_models_provider_model_name ON modelkit_models (provider, model_name);
