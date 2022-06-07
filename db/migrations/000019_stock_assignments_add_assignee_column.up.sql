BEGIN;

ALTER TABLE stock_assignments
ADD COLUMN assignee_id int;

ALTER TABLE stock_assignments
ADD CONSTRAINT fk_assignee_id
FOREIGN KEY (assignee_id)
REFERENCES users(id);

COMMIT;
