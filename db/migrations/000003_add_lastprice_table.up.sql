BEGIN;
  CREATE TABLE IF NOT EXISTS lastprices (
    updatedat TIMESTAMP,
    coinid NUMERIC UNIQUE,
    price NUMERIC,
    CONSTRAINT coins_id_fkey FOREIGN KEY (coinid)
      REFERENCES coins (coinid) ON DELETE CASCADE
  );
COMMIT;
