BEGIN;

CREATE TYPE import_status AS ENUM (
	'pending',
	'imported',
	'failed'
);

CREATE TABLE IF NOT EXISTS procurements (
	id BIGSERIAL PRIMARY KEY,
	filename TEXT NOT NULL,
	status import_status DEFAULT 'pending',
	failed_reason TEXT NULL,

	created_at timestamp NOT NULL DEFAULT NOW(),
	updated_at timestamp NULL DEFAULT current_timestamp,
	deleted_at timestamp
);

COMMIT;
