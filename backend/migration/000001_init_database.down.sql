-- 删除外键约束
ALTER TABLE model_provider_models DROP CONSTRAINT IF EXISTS fk_model_provider_models_provider_id;

-- 删除索引
DROP INDEX IF EXISTS idx_model_name_api_base_type;
DROP INDEX IF EXISTS unique_idx_model_provider_models_provider_id_name;
DROP INDEX IF EXISTS unique_idx_model_providers_name;

-- 删除表
DROP TABLE IF EXISTS models CASCADE;
DROP TABLE IF EXISTS model_provider_models CASCADE;
DROP TABLE IF EXISTS model_providers CASCADE;

-- 删除扩展
DROP EXTENSION IF EXISTS "uuid-ossp";