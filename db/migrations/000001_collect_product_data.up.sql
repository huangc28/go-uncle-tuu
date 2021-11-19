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

