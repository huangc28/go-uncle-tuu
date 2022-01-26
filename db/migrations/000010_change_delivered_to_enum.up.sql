BEGIN;

CREATE TYPE delivered_status AS ENUM (
    'not_yet_reported',
    'delivered',
    'not_delivered'
);

ALTER TABLE inventory
ADD COLUMN delivered delivered_status DEFAULT 'not_yet_reported';

COMMIT;
