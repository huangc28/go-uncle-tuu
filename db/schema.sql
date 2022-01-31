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
BEGIN;
ALTER TABLE users
ADD COLUMN uuid VARCHAR(20) UNIQUE NOT NULL;
COMMIT;
BEGIN;
ALTER TABLE inventory
ADD COLUMN IF NOT EXISTS reserved_for_user int;

ALTER TABLE inventory
   ADD CONSTRAINT fk_reserve_for_user_id
   FOREIGN KEY (reserved_for_user)
   REFERENCES users(id);

COMMIT;
BEGIN;
ALTER TABLE inventory
DROP COLUMN IF EXISTS delivered;
COMMIT;
BEGIN;

CREATE TYPE delivered_status AS ENUM (
    'not_yet_reported',
    'delivered',
    'not_delivered'
);

ALTER TABLE inventory
ADD COLUMN delivered delivered_status DEFAULT 'not_yet_reported';

COMMIT;
BEGIN;

ALTER TABLE inventory
ADD COLUMN temp_receipt text;

COMMIT;
