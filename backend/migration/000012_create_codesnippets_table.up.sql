-- Create code_snippets table
CREATE TABLE IF NOT EXISTS code_snippets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_file_id UUID NOT NULL,
    -- Basic Info
    name VARCHAR(255) NOT NULL,
    snippet_type VARCHAR(255) NOT NULL,
    language VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    hash VARCHAR(255) NOT NULL,
    -- Position Info
    start_line INTEGER NOT NULL,
    end_line INTEGER NOT NULL,
    start_column INTEGER NOT NULL,
    end_column INTEGER NOT NULL,
    -- Context Info
    namespace VARCHAR(255),
    container_name VARCHAR(255),
    scope JSONB,
    dependencies JSONB,
    -- Structured Info
    parameters JSONB,
    signature TEXT,
    definition_text TEXT,
    structured_info JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- Create indexes for code_snippets table
-- Hash index for duplicate detection and fast lookup
CREATE INDEX IF NOT EXISTS codesnippet_hash ON code_snippets (hash);
-- Workspace file ID index for querying snippets by file
CREATE INDEX IF NOT EXISTS codesnippet_workspace_file_id ON code_snippets (workspace_file_id);
-- Language and snippet type index for filtering
CREATE INDEX IF NOT EXISTS codesnippet_language_type ON code_snippets (language, snippet_type);
-- Language and name index for searching
CREATE INDEX IF NOT EXISTS codesnippet_language_name ON code_snippets (language, name);
-- Add foreign key constraint
ALTER TABLE code_snippets
ADD CONSTRAINT fk_codesnippets_workspace_file_id FOREIGN KEY (workspace_file_id) REFERENCES workspace_files(id) ON DELETE CASCADE;