BEGIN;

ALTER TABLE users
ADD COLUMN can_be_exported BOOLEAN DEFAULT true;

COMMIT;
