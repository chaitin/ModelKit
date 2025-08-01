CREATE TABLE extensions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
    version VARCHAR(255) NOT NULL,
    path VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX unique_idx_extensions_version ON extensions(version);