BEGIN;
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
COMMIT;
