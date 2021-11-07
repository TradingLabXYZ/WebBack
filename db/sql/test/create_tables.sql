CREATE TYPE privacies AS ENUM ('all', 'private', 'subscribers', 'followers');
CREATE TYPE plans AS ENUM ('basic', 'premium', 'pro');

CREATE TABLE IF NOT EXISTS users (
  wallet VARCHAR(42) NOT NULL UNIQUE,
  username VARCHAR(255) NOT NULL UNIQUE,
  privacy privacies NOT NULL,
  profilepicture TEXT,
  twitter TEXT,
  website TEXT,
  plan plans NOT NULL,
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sessions (
  code VARCHAR(64) NOT NULL UNIQUE,
  userwallet VARCHAR(42) NOT NULL,
  createdat TIMESTAMP NOT NULL,
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

CREATE TABLE IF NOT EXISTS trades (
  code VARCHAR(12) NOT NULL UNIQUE,
  userwallet VARCHAR(42) NOT NULL,
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  exchange VARCHAR(64),
  firstpair NUMERIC NOT NULL REFERENCES coins(coinid),
  secondpair NUMERIC NOT NULL REFERENCES coins(coinid),
  isopen BOOLEAN,
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
