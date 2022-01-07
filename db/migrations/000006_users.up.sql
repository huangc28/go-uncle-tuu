BEGIN;

CREATE TYPE roles AS ENUM (
	'admin',
	'vendor',
	'client'
);

CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,

	username VARCHAR UNIQUE NOT NULL,
	password VARCHAR NOT NULL,
	role roles DEFAULT 'admin',

	created_at timestamp NOT NULL DEFAULT NOW(),
	updated_at timestamp NULL DEFAULT current_timestamp,
	deleted_at timestamp
);

COMMIT;
