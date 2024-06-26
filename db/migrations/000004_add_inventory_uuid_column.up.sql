BEGIN;

CREATE extension IF NOT EXISTS "uuid-ossp";

COMMIT;

BEGIN;

ALTER TABLE inventory
ADD COLUMN uuid VARCHAR(40) UNIQUE NOT NULL DEFAULT uuid_generate_v1();

COMMIT;
