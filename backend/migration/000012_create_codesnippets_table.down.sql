-- Drop foreign key constraint
ALTER TABLE code_snippets DROP CONSTRAINT IF EXISTS fk_codesnippets_workspace_file_id;
-- Drop indexes for code_snippets table
DROP INDEX IF EXISTS codesnippet_hash;
DROP INDEX IF EXISTS codesnippet_workspace_file_id;
DROP INDEX IF EXISTS codesnippet_language_type;
DROP INDEX IF EXISTS codesnippet_language_name;
-- Drop code_snippets table
DROP TABLE IF EXISTS code_snippets;