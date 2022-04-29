BEGIN;
  ALTER TABLE sessions
  DROP COLUMN timezone;
COMMIT;
