ALTER TABLE invite_codes ADD COLUMN status VARCHAR(255);
UPDATE invite_codes SET status = 'used' WHERE status IS NULL;
ALTER TABLE invite_codes ADD COLUMN expired_at TIMESTAMPTZ;