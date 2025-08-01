-- 添加prompt字段（如果不存在）
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS prompt text;

-- 修改cursor_position为JSON类型（如果还是bigint类型）
ALTER TABLE tasks ALTER COLUMN cursor_position TYPE jsonb USING
  CASE
    WHEN cursor_position IS NOT NULL THEN
      jsonb_build_object('position', cursor_position, 'line', 1, 'column', cursor_position)
    ELSE NULL
  END;
