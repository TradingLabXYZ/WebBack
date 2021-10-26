CREATE TYPE privacies AS ENUM ('all', 'private', 'subscribers', 'followers');
CREATE TYPE plans AS ENUM ('basic', 'premium', 'pro');

CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  code VARCHAR(12) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  username VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  privacy privacies NOT NULL,
  profilepicture TEXT,
  twitter TEXT,
  website TEXT,
  plan plans NOT NULL,
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  deletedat TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL,
  userid INTEGER REFERENCES users(id),
  createdat TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS coins (
  id SERIAL PRIMARY KEY,
  coinid NUMERIC UNIQUE,
  name TEXT,
  symbol TEXT,
  slug TEXT
);

CREATE TABLE IF NOT EXISTS prices (
  id SERIAL PRIMARY KEY,
  createdat TIMESTAMP,
  coinid NUMERIC REFERENCES coins(coinid),
  price NUMERIC
);

CREATE TABLE IF NOT EXISTS trades (
  code VARCHAR(12) NOT NULL UNIQUE,
  usercode VARCHAR(12) REFERENCES users(code),
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  exchange VARCHAR(64),
  firstpair NUMERIC NOT NULL REFERENCES coins(coinid),
  secondpair NUMERIC NOT NULL REFERENCES coins(coinid),
  isopen BOOLEAN
);

CREATE TABLE IF NOT EXISTS subtrades (
  code VARCHAR(12) NOT NULL UNIQUE,
  tradecode VARCHAR(12) NOT NULL,
  usercode VARCHAR(12) NOT NULL REFERENCES users(code),
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  type VARCHAR(5),
  reason VARCHAR(64),
  quantity NUMERIC,
  avgprice NUMERIC,
  total NUMERIC,
  CONSTRAINT trades_code_fkey FOREIGN KEY (tradecode)
    REFERENCES trades (code) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS followers (
  id SERIAL PRIMARY KEY,
  followefrom INTEGER REFERENCES users(id) NOT NULL,
  followto INTEGER references users(id) NOT NULL,
  createdat TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscribers (
  id SERIAL PRIMARY KEY,
  subscribefrom INTEGER REFERENCES users(id) NOT NULL,
  subscribeto INTEGER REFERENCES users(id) NOT NULL,
  createdat TIMESTAMP
);

CREATE TABLE IF NOT EXISTS internalwallets (
  id SERIAL PRIMARY KEY,
  blockchain TEXT,
  currency TEXT,
  address TEXT,
  description TEXT,
  createdat TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS memos (
  id SERIAL PRIMARY KEY,
  userid INTEGER REFERENCES users(id) NOT NULL,
  blockchain TEXT NOT NULL,
  currency TEXT NOT NULL,
  depositaddress TEXT NOT NULL,
  memo VARCHAR(20) NOT NULL,
  createdat TIMESTAMP NOT NULL,
  status TEXT
);

CREATE TABLE IF NOT EXISTS payments (
  id SERIAL PRIMARY KEY,
  userid INTEGER REFERENCES users(id) NOT NULL,
  type TEXT NOT NULL,
  blockchain TEXT NOT NULL,
  currency TEXT NOT NULL,
  transactionid TEXT NOT NULL,
  amount NUMERIC NOT NULL,
  months INTEGER NOT NULL,
  createdat TIMESTAMP NOT NULL,
  endat TIMESTAMP NOT NULL
);

CREATE OR REPLACE FUNCTION notify_changes()
  RETURNS trigger AS $$
  BEGIN
    PERFORM pg_notify(
      'activity_update',
      OLD.usercode
    );
    RETURN OLD;
  END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS activity_update ON trades;
CREATE TRIGGER activity_update
AFTER INSERT OR UPDATE OR DELETE
ON trades
FOR EACH ROW
EXECUTE PROCEDURE notify_changes();

DROP TRIGGER IF EXISTS activity_update ON subtrades;
CREATE TRIGGER activity_update
AFTER INSERT OR UPDATE OR DELETE
ON subtrades
FOR EACH ROW
EXECUTE PROCEDURE notify_changes();
