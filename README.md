## Add game item to inventory

```postgres
BEGIN;

LOCK TABLE inventory IN ROW EXCLUSIVE MODE;

SELECT * AS target_prod
FROM inventory
INNER JOIN product_info ON inventory.prod_id = product_info.id
WHERE PRODUCT_INFO = $1 FOR UPDATE OF inventory;

UPDATE inventory
SET quantity = quantity + 1
FROM product_info
WHERE
	inventory.prod_id = product_info.id AND
	product_info.prod_id = $1;
COMMIT;
```
