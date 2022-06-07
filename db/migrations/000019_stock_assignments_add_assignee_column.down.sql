BEGIN;

ALTER TABLE stock_assignments
DROP CONSTRAINT fk_assignee_id;

ALTER TABLE stock_assignments
DROP COLUMN assignee_id;

COMMIT;
