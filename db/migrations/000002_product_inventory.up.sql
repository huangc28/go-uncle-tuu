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
