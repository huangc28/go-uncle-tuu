BEGIN;

ALTER TABLE inventory
DROP COLUMN IF EXISTS temp_receipt;

COMMIT;