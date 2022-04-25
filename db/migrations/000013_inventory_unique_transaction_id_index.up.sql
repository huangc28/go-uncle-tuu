BEGIN;

CREATE UNIQUE INDEX transaction_id_idx on inventory(transaction_id);

COMMIT;
