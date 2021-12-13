BEGIN;

ALTER TABLE inventory
ADD COLUMN delivered boolean NOT NULL DEFAULT false;

COMMIT;
