BEGIN;
  CREATE TABLE IF NOT EXISTS smartcontractevents (
    createdat TIMESTAMP NOT NULL,
    transaction VARCHAR(66) NOT NULL,
    contract VARCHAR(42) NOT NULL,
    sender VARCHAR(42) NOT NULL,
    name TEXT NOT NULL,
    signature VARCHAR(66) NULL,
    payload JSON NOT NULL
  );
COMMIT;
