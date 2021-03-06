CREATE TYPE privacies AS ENUM ('all', 'private', 'subscribers', 'followers');

CREATE TABLE IF NOT EXISTS users (
  wallet VARCHAR(42) NOT NULL UNIQUE,
  username VARCHAR(20),
  twitter VARCHAR(15),
  discord VARCHAR(30),
  github VARCHAR(32),
  privacy privacies NOT NULL,
  profilepicture TEXT,
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  code VARCHAR(64) NOT NULL UNIQUE,
  userwallet VARCHAR(42) NOT NULL,
  createdat TIMESTAMP NOT NULL,
  origin TEXT NOT NULL,
  timezone TEXT NOT NULL,
  CONSTRAINT users_userwallet_fkey FOREIGN KEY (userwallet)
    REFERENCES users (wallet) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS coins (
  coinid NUMERIC UNIQUE,
  name TEXT,
  symbol TEXT,
  slug TEXT
);

CREATE TABLE IF NOT EXISTS prices (
  createdat TIMESTAMP,
  coinid NUMERIC,
  price NUMERIC,
  CONSTRAINT coins_id_fkey FOREIGN KEY (coinid)
    REFERENCES coins (coinid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS lastprices (
  updatedat TIMESTAMP,
  coinid NUMERIC UNIQUE,
  price NUMERIC,
  CONSTRAINT coins_id_fkey FOREIGN KEY (coinid)
    REFERENCES coins (coinid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS trades (
  code VARCHAR(12) NOT NULL UNIQUE,
  userwallet VARCHAR(42) NOT NULL,
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  exchange VARCHAR(64),
  firstpair NUMERIC NOT NULL REFERENCES coins(coinid),
  secondpair NUMERIC NOT NULL REFERENCES coins(coinid),
  CONSTRAINT users_wallet_fkey FOREIGN KEY (userwallet)
    REFERENCES users (wallet) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subtrades (
  code VARCHAR(12) NOT NULL UNIQUE,
  tradecode VARCHAR(12) NOT NULL,
  userwallet VARCHAR(42) NOT NULL,
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  type VARCHAR(5),
  reason VARCHAR(64),
  quantity NUMERIC,
  avgprice NUMERIC,
  total NUMERIC,
  CONSTRAINT users_wallet_fkey FOREIGN KEY (userwallet)
    REFERENCES users (wallet) ON DELETE CASCADE,
  CONSTRAINT trades_code_fkey FOREIGN KEY (tradecode)
    REFERENCES trades (code) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS followers (
  followfrom VARCHAR(42) NOT NULL,
  followto VARCHAR(42) NOT NULL,
  createdat TIMESTAMP,
  CONSTRAINT users_userwallet_followfrom_fkey FOREIGN KEY (followfrom)
    REFERENCES users (wallet) ON DELETE CASCADE,
  CONSTRAINT users_userwallet_followto_fkey FOREIGN KEY (followto)
    REFERENCES users (wallet) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS subscribers (
  subscribefrom VARCHAR(42) NOT NULL,
  subscribeto VARCHAR(42) NOT NULL,
  createdat TIMESTAMP,
  CONSTRAINT users_userwallet_subscribefrom_fkey FOREIGN KEY (subscribefrom)
    REFERENCES users (wallet) ON DELETE CASCADE,
  CONSTRAINT users_userwallet_subscribeto_fkey FOREIGN KEY (subscribeto)
    REFERENCES users (wallet) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS smartcontractevents (
  createdat TIMESTAMP NOT NULL,
  transaction VARCHAR(66) NOT NULL,
  contract VARCHAR(42) NOT NULL,
  sender VARCHAR(42) NOT NULL,
  name TEXT NOT NULL,
  signature VARCHAR(66) NULL,
  payload JSON NOT NULL
);

CREATE TABLE IF NOT EXISTS visibilities (
  wallet VARCHAR(42) NOT NULL UNIQUE,
  totalcounttrades BOOLEAN NOT NULL,
  totalportfolio BOOLEAN NOT NULL,  
  totalreturn BOOLEAN NOT NULL,
  totalroi BOOLEAN NOT NULL,
  tradeqtyavailable BOOLEAN NOT NULL,
  tradevalue BOOLEAN NOT NULL,
  tradereturn BOOLEAN NOT NULL,
  traderoi BOOLEAN NOT NULL,
  subtradesall BOOLEAN NOT NULL,
  subtradereasons BOOLEAN NOT NULL,
  subtradequantity BOOLEAN NOT NULL,
  subtradeavgprice BOOLEAN NOT NULL,
  subtradetotal BOOLEAN NOT NULL,
  CONSTRAINT users_userwallet_fkey FOREIGN KEY (wallet)
    REFERENCES users (wallet) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS competitions (
  name TEXT NOT NULL UNIQUE,
  submissionstartedat TIMESTAMP NOT NULL,
  submissionendedat TIMESTAMP NOT NULL,
  competitionstartedat TIMESTAMP NOT NULL,
  competitionendedat TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS submissions (
  updatedat TIMESTAMP NOT NULL,
  competitionname TEXT NOT NULL,
  userwallet VARCHAR(42) NOT NULL,
  payload JSON NOT NULL,
  CONSTRAINT users_wallet_fkey FOREIGN KEY (userwallet)
    REFERENCES users (wallet) ON DELETE CASCADE,
  CONSTRAINT competitions_name_fkey FOREIGN KEY (competitionname)
    REFERENCES competitions (name) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION notify_changes()
  RETURNS trigger AS $$
  BEGIN
    PERFORM pg_notify(
      'activity_update',
      CASE WHEN OLD.userwallet IS NULL THEN NEW.userwallet ELSE OLD.userwallet END
    );
    RETURN OLD;
  END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS activity_update ON subtrades;
CREATE TRIGGER activity_update
AFTER INSERT OR UPDATE OR DELETE
ON subtrades
FOR EACH ROW
EXECUTE PROCEDURE notify_changes();
