BEGIN;

CREATE TABLE IF NOT EXISTS product_info (
	id BIGSERIAL PRIMARY KEY,
	prod_name TEXT,
	prod_id TEXT NOT NULL,
	prod_desc TEXT,
	price numeric(12, 2),
	game_bundle_id TEXT NOT NULL
);

CREATE UNIQUE INDEX unique_product_data
ON product_info(prod_id, game_bundle_id);

COMMIT;

BEGIN;

CREATE TABLE IF NOT EXISTS inventory (
	id BIGSERIAL PRIMARY KEY,
	prod_id INT REFERENCES product_info (id),
	transaction_id text,
	receipt text,
	available BOOLEAN DEFAULT true,

	created_at timestamp NOT NULL DEFAULT NOW(),
	updated_at timestamp NULL DEFAULT current_timestamp,
	deleted_at timestamp
);

COMMIT;
BEGIN;

ALTER TABLE inventory
ADD column transaction_time timestamp NOT NULL DEFAULT NOW();

COMMIT;
BEGIN;

CREATE extension IF NOT EXISTS "uuid-ossp";

COMMIT;

BEGIN;

ALTER TABLE inventory
ADD COLUMN uuid VARCHAR(40) UNIQUE NOT NULL DEFAULT uuid_generate_v1();

COMMIT;
BEGIN;

ALTER TABLE inventory
ADD COLUMN delivered boolean NOT NULL DEFAULT false;

COMMIT;
