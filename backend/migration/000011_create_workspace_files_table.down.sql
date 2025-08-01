-- Drop workspace_files table first (due to foreign key dependency)
DROP TABLE IF EXISTS workspace_files CASCADE;

-- Drop workspaces table and all its constraints and indexes
DROP TABLE IF EXISTS workspaces CASCADE;
