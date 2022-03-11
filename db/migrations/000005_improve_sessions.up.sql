BEGIN;
  ALTER TABLE sessions ADD COLUMN origin TEXT; 
  UPDATE sessions SET origin = 'web' WHERE 1 = 1;
  ALTER TABLE sessions ALTER COLUMN origin SET NOT NULL;
COMMIT;
