-- Create workspaces table
CREATE TABLE IF NOT EXISTS workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255),
    description TEXT,
    root_path VARCHAR(255) NOT NULL,
    settings JSONB,
    last_accessed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for workspaces table
-- Unique constraint: user_id + root_path (same user cannot have duplicate root paths)
CREATE UNIQUE INDEX IF NOT EXISTS workspaces_user_id_root_path ON workspaces (user_id, root_path);

-- User ID index
CREATE INDEX IF NOT EXISTS workspaces_user_id ON workspaces (user_id);

-- Last accessed time index
CREATE INDEX IF NOT EXISTS workspaces_last_accessed_at ON workspaces (last_accessed_at);

-- Name index for searching by name
CREATE INDEX IF NOT EXISTS workspaces_name ON workspaces (name);

-- Root path index
CREATE INDEX IF NOT EXISTS workspaces_root_path ON workspaces (root_path);

-- Create workspace_files table
CREATE TABLE IF NOT EXISTS workspace_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    workspace_id UUID NOT NULL,
    path VARCHAR(255) NOT NULL,
    content TEXT,
    hash VARCHAR(255) NOT NULL,
    language VARCHAR(255),
    size BIGINT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for workspace_files table
-- Unique constraint: user_id + workspace_id + path (same user/workspace cannot have duplicate file paths)
CREATE UNIQUE INDEX IF NOT EXISTS workspacefile_user_id_workspace_id_path ON workspace_files (user_id, workspace_id, path);

-- Hash index for duplicate detection and fast lookup
CREATE INDEX IF NOT EXISTS workspacefile_hash ON workspace_files (hash);

-- Workspace + hash index to prevent duplicate content in same workspace
CREATE INDEX IF NOT EXISTS workspacefile_workspace_id_hash ON workspace_files (workspace_id, hash);

-- Language index for filtering by programming language
CREATE INDEX IF NOT EXISTS workspacefile_language ON workspace_files (language);

-- Updated time index for finding recently modified files
CREATE INDEX IF NOT EXISTS workspacefile_updated_at ON workspace_files (updated_at);

-- File size index for statistics and size-based queries
CREATE INDEX IF NOT EXISTS workspacefile_size ON workspace_files (size);

-- Workspace ID index for querying all files in a workspace
CREATE INDEX IF NOT EXISTS workspacefile_workspace_id ON workspace_files (workspace_id);

-- User ID index for querying all files by user
CREATE INDEX IF NOT EXISTS workspacefile_user_id ON workspace_files (user_id);

-- Add foreign key constraints
ALTER TABLE workspaces ADD CONSTRAINT fk_workspaces_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE workspace_files ADD CONSTRAINT fk_workspace_files_workspace_id 
    FOREIGN KEY (workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE;

ALTER TABLE workspace_files ADD CONSTRAINT fk_workspace_files_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
