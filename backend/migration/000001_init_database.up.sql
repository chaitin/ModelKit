-- 创建 UUID 扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 创建模型表
CREATE TABLE IF NOT EXISTS modelkit_models (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    model_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    api_base VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    api_version VARCHAR(255),
    api_header TEXT,
    provider VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_modelkit_models_name_api_base_type ON modelkit_models (model_name, api_base, model_type);
