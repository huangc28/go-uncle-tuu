BEGIN;

ALTER TABLE inventory
ADD column transaction_time timestamp NOT NULL DEFAULT NOW();

COMMIT;
