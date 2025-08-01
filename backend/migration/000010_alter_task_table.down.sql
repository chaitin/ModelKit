-- 回滚cursor_position类型修改
ALTER TABLE tasks ALTER COLUMN cursor_position TYPE bigint USING
  CASE
    WHEN cursor_position IS NOT NULL AND cursor_position ? 'position' THEN
      (cursor_position->>'position')::bigint
    ELSE NULL
  END;

-- 删除prompt字段
ALTER TABLE tasks DROP COLUMN IF EXISTS prompt;
