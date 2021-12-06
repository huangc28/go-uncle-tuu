BEGIN;

ALTER TABLE inventory
DROP COLUMN transaction_time,

COMMIT;
