BEGIN;

DROP TABLE IF EXISTS stock_assignments;

ALTER TABLE inventory
DROP CONSTRAINT fk_assignment_id;

ALTER TABLE inventory
DROP COLUMN assignment_id;

COMMIT;
