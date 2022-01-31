BEGIN;

ALTER TABLE inventory
ADD COLUMN temp_receipt text;

COMMIT;
