BEGIN;
  ALTER TABLE sessions
  ADD COLUMN timezone TEXT;
COMMIT;

BEGIN;
  UPDATE sessions
  SET timezone = 'Europe/Berlin'
  WHERE 1 = 1;
COMMIT;


BEGIN;
  ALTER TABLE sessions
  ALTER COLUMN timezone SET NOT NULL;
COMMIT;
