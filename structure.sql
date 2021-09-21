create table users (
  id            serial primary key,
  email         varchar(255) not null unique,
  username      varchar(255) not null unique,
  password      varchar(255) not null,
  createdat     timestamp not null,
  updatedat     timestamp not null,
  deletedat     timestamp
);

create table sessions (
  id            serial primary key,
  uuid          varchar(64) not null unique,
  email         varchar(255),
  userid        integer references users(id),
  createdat     timestamp not null,
  deletedat     timestamp
);

CREATE TABLE coins (
  id serial primary key,
  coinid numeric UNIQUE,
  name text,
  symbol text,
  slug text
);

CREATE TABLE prices (
  id serial primary key,
  createdat timestamp,
  coinid numeric references coins(coinid),
  price numeric
);

create table trades (
  id              serial primary key,
  userid          integer references users(id),
  usertrade       integer not null,
  createdat       timestamp not null,
  updatedat       timestamp not null,
  deletedat       timestamp,
  exchange        varchar(64),
  firstpair       numeric references coins(coinid),
  secondpair      numeric references coins(coinid),
  isopen          boolean
);

create table subtrades (
  subtradeid      integer,
  tradeid         integer references trades(id),
  createdat       timestamp not null,
  updatedat       timestamp not null,
  deletedat       timestamp,
  tradetimestamp  timestamp,
  type            varchar(5),
  reason          varchar(64),
  quantity        numeric,
  avgprice        numeric,
  total           numeric
);
