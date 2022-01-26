BEGIN;

ALTER TABLE inventory
ALTER ADD COLUMN delivered boolean;

COMMIT;
