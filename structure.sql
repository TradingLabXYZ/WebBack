CREATE TYPE privacies AS ENUM ('all', 'private', 'subscribers', 'followers');
CREATE TYPE plans AS ENUM ('basic', 'premium', 'pro');

CREATE TABLE users (
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

CREATE TABLE sessions (
  id SERIAL PRIMARY KEY,
  uuid VARCHAR(64) NOT NULL UNIQUE,
  email VARCHAR(255),
  userid INTEGER REFERENCES users(id),
  createdat TIMESTAMP NOT NULL,
  deletedat TIMESTAMP
);

CREATE TABLE coins (
  id SERIAL PRIMARY KEY,
  coinid NUMERIC UNIQUE,
  name TEXT,
  symbol TEXT,
  slug TEXT
);

CREATE TABLE prices (
  id SERIAL PRIMARY KEY,
  createdat TIMESTAMP,
  coinid NUMERIC REFERENCES coins(coinid),
  price NUMERIC
);

CREATE TABLE trades (
  id VARCHAR(12) NOT NULL UNIQUE,
  userid INTEGER REFERENCES users(id),
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  deletedat TIMESTAMP,
  exchange VARCHAR(64),
  firstpair NUMERIC REFERENCES coins(coinid),
  secondpair NUMERIC REFERENCES coins(coinid),
  isopen BOOLEAN
);

CREATE TABLE subtrades (
  id SERIAL PRIMARY KEY,
  tradeid VARCHAR(12) REFERENCES trades(id),
  createdat TIMESTAMP NOT NULL,
  updatedat TIMESTAMP NOT NULL,
  deletedat TIMESTAMP,
  tradetimestamp TIMESTAMP,
  type VARCHAR(5),
  reason VARCHAR(64),
  quantity NUMERIC,
  avgprice NUMERIC,
  total NUMERIC
);

CREATE TABLE followers (
  id SERIAL PRIMARY KEY,
  followefrom INTEGER REFERENCES users(id) NOT NULL,
  followto INTEGER references users(id) NOT NULL,
  createdat TIMESTAMP
);

CREATE TABLE subscribers (
  id SERIAL PRIMARY KEY,
  subscribefrom INTEGER REFERENCES users(id) NOT NULL,
  subscribeto INTEGER REFERENCES users(id) NOT NULL,
  createdat TIMESTAMP
);

CREATE TABLE internalwallets (
  id SERIAL PRIMARY KEY,
  blockchain TEXT,
  currency TEXT,
  address TEXT,
  description TEXT,
  createdat TIMESTAMP NOT NULL
);

CREATE TABLE memos (
  id SERIAL PRIMARY KEY,
  userid INTEGER REFERENCES users(id) NOT NULL,
  blockchain TEXT NOT NULL,
  currency TEXT NOT NULL,
  depositaddress TEXT NOT NULL,
  memo VARCHAR(20) NOT NULL,
  createdat TIMESTAMP NOT NULL,
  status BOOLEAN 
);

CREATE TABLE payments (
  id SERIAL PRIMARY KEY,
  userid INTEGER REFERENCES users(id) NOT NULL,
  blockchain TEXT NOT NULL,
  currency TEXT NOT NULL,
  transactionid TEXT NOT NULL,
  amount NUMERIC NOT NULL,
  months INTEGER NOT NULL,
  createdat TIMESTAMP NOT NULL,
  endat TIMESTAMP NOT NULL
);
