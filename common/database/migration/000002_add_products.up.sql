CREATE TABLE products (
    product_id TEXT PRIMARY KEY NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    quantity INT,
    price INT,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL
);