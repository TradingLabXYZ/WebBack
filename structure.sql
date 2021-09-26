create table users (
  id            serial primary key,
  email         varchar(255) not null unique,
  username      varchar(255) not null unique,
  password      varchar(255) not null,
  permission    permissions not null,
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

CREATE TABLE trades (
  id              varchar(12) not null unique,
  userid          integer references users(id),
  createdat       timestamp not null,
  updatedat       timestamp not null,
  deletedat       timestamp,
  exchange        varchar(64),
  firstpair       numeric references coins(coinid),
  secondpair      numeric references coins(coinid),
  isopen          boolean
);

CREATE TABLE subtrades (
  id              serial primary key,
  tradeid         varchar(12) references trades(id),
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

CREATE TABLE followers (
  id serial primary key,
  usera integer references users(id) not null,
  userb integer references users(id) not null,
  createdat timestamp,
  updatedat timestamp,
  deletedat timestamp
);

CREATE TABLE subscribers (
  id serial primary key,
  usera integer references users(id) not null,
  userb integer references users(id) not null,
  createdat timestamp,
  updatedat timestamp,
  deletedat timestamp
);

CREATE TABLE individuals (
  id serial primary key,
  usera integer references users(id) not null,
  userb integer references users(id) not null,
  createdat timestamp,
  updatedat timestamp,
  deletedat timestamp
);
