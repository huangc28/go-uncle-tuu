BEGIN;
ALTER TABLE users
ADD COLUMN uuid VARCHAR(20) UNIQUE NOT NULL;
COMMIT;
