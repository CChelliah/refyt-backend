ALTER TABLE products
RENAME COLUMN title TO product_name;

ALTER TABLE products
ADD COLUMN rrp INTEGER,
ADD COLUMN fit_notes TEXT,
ADD COLUMN designer TEXT;
