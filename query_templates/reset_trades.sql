DROP TABLE subtrades;
DROP TABLE trades;

create table trades (
  id              serial primary key,
  userid          integer references users(id),
  usertrade       integer not null,
  createdat       timestamp not null,
  updatedat       timestamp not null,
  deletedat       timestamp,
  exchange        varchar(64),
  firstpair       varchar(20),
  secondpair      varchar(20),
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

INSERT INTO trades (userid, usertrade, createdat, updatedat, exchange, firstpair, secondpair, isopen) VALUES
(8, 1, '2021-09-13 05:46:44.173258', '2021-09-13 05:46:44.173258', 'Bittrex', 'BTC', 'LUNA', TRUE);

INSERT INTO subtrades (subtradeid, tradeid, createdat, updatedat, tradetimestamp, type, reason, quantity, avgprice, total)  VALUES
(1, 1, '2021-09-13 05:46:44.173258', '2021-09-13 05:46:44.173258', '2021-09-10 02:09:00', 'BUY', 'Suggestion', 9.59126633, 0.00082135, 0.00793687),
(2, 1, '2021-09-13 05:46:44.224704', '2021-09-13 05:46:44.224704', '2021-09-12 11:09:00', 'SELL', 'Profit', 4.79563317, 0.00087013, 0.00414153),
(3, 1, '2021-09-13 05:46:44.272248', '2021-09-13 05:46:44.272248', '2021-09-12 11:09:00', 'BUY', 'Local Min', 2.54308536, 0.00084389, 0.00216217),
(4, 1, '2021-09-13 05:46:44.324024', '2021-09-13 05:46:44.324024', '2021-09-12 11:09:00', 'SELL', 'Profit', 3.66935926, 0.00094858, 0.00345458),
(5, 1, '2021-09-13 05:46:44.373735', '2021-09-13 05:46:44.373735', '2021-09-12 02:09:00', 'BUY', 'Local Min', 3.0606912, 0.00088109, 0.00271697),
(6, 1, '2021-09-13 05:46:44.425591', '2021-09-13 05:46:44.425591', '2021-09-12 10:09:00', 'BUY', 'Local Min', 1.57605392, 0.00085554, 0.00135849),
(7, 1, '2021-09-13 05:46:44.477528', '2021-09-13 05:46:44.477528', '2021-09-13 07:46:00', 'BUY', 'Local Min', 1.62911214, 0.00082767, 0.00135848);

INSERT INTO trades (userid, usertrade, createdat, updatedat, exchange, firstpair, secondpair, isopen) VALUES
(8, 2, '2021-09-13 05:46:44.173258', '2021-09-13 05:46:44.173258', 'Bittrex', 'BTC', 'DOT', FALSE);

INSERT INTO subtrades (subtradeid, tradeid, createdat, updatedat, tradetimestamp, type, reason, quantity, avgprice, total)  VALUES
(1, 2, '2021-09-13 05:46:44.173258', '2021-09-13 05:46:44.173258', '2021-09-10 02:09:00', 'BUY', 'Suggestion', 9.59126633, 0.00082135, 0.00793687),
(2, 2, '2021-09-13 05:46:44.224704', '2021-09-13 05:46:44.224704', '2021-09-12 11:09:00', 'SELL', 'Profit', 9.59126633, 0.00082135, 0.0078111);
