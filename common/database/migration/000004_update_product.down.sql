ALTER TABLE products
RENAME COLUMN product_name TO title;

ALTER TABLE products
DROP COLUMN rrp,
DROP COLUMN  fit_notes,
DROP COLUMN  designer;