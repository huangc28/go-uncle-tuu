BEGIN;

ALTER TABLE inventory
DROP COLUMN IF EXISTS reserved_for_user;

COMMIT;
