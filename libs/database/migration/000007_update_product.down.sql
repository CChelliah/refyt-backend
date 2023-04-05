ALTER TABLE products
DROP COLUMN category,
DROP COLUMN shipping_price,
DROP COLUMN size,
DROP COLUMN images,
ADD COLUMN quantity INTEGER;