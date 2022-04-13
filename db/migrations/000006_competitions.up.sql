BEGIN;
  CREATE TABLE IF NOT EXISTS competitions (
    name TEXT NOT NULL UNIQUE,
    submissionstartedat TIMESTAMP NOT NULL,
    submissionendedat TIMESTAMP NOT NULL,
    competitionstartedat TIMESTAMP NOT NULL,
    competitionendedat TIMESTAMP NOT NULL
  );
COMMIT;
BEGIN;
  CREATE TABLE IF NOT EXISTS submissions (
    competitionname TEXT NOT NULL,
    userwallet VARCHAR(42) NOT NULL,
    payload JSON NOT NULL,
    CONSTRAINT users_wallet_fkey FOREIGN KEY (userwallet)
      REFERENCES users (wallet) ON DELETE CASCADE,
    CONSTRAINT competitions_name_fkey FOREIGN KEY (competitionname)
      REFERENCES competitions (name) ON DELETE CASCADE
  );
COMMIT;
