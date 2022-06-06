BEGIN;

CREATE TABLE IF NOT EXISTS stock_assignments (
	id BIGSERIAL PRIMARY KEY,
	uuid VARCHAR(40) UNIQUE NOT NULL DEFAULT uuid_generate_v1(),

	created_at timestamp NOT NULL DEFAULT NOW(),
	updated_at timestamp NULL DEFAULT current_timestamp,
	deleted_at timestamp
);

ALTER TABLE inventory
ADD COLUMN IF NOT EXISTS assignment_id int;

ALTER TABLE inventory
ADD CONSTRAINT fk_assignment_id
FOREIGN KEY (assignment_id)
REFERENCES stock_assignments(id);

COMMIT;
