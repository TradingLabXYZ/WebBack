BEGIN;

  DROP TABLE IF EXISTS prices;
  DROP TABLE IF EXISTS subtrades;
  DROP TABLE IF EXISTS trades;
  DROP TABLE IF EXISTS followers;
  DROP TABLE IF EXISTS subscribers;
  DROP TABLE IF EXISTS internalwallets;
  DROP TABLE IF EXISTS memos;
  DROP TABLE IF EXISTS payments;
  DROP TABLE IF EXISTS sessions;
  DROP TABLE IF EXISTS users;
  DROP TABLE IF EXISTS coins;
  DROP TYPE IF EXISTS privacies;

COMMIT;
