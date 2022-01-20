BEGIN;
ALTER TABLE inventory
ADD COLUMN IF NOT EXISTS reserved_for_user int;

ALTER TABLE inventory
   ADD CONSTRAINT fk_reserve_for_user_id
   FOREIGN KEY (reserved_for_user)
   REFERENCES users(id);

COMMIT;
