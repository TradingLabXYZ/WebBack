BEGIN;
  CREATE TABLE IF NOT EXISTS contractplans (
    createdat TIMESTAMP NOT NULL,
    transaction VARCHAR(66) NOT NULL,
    sender VARCHAR(42) NOT NULL,
    contract VARCHAR(42) NOT NULL,
    value text NOT NULL
  );
COMMIT;
