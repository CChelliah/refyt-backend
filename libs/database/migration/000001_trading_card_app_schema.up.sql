CREATE TABLE users (
    uid TEXT NOT NULL,
    email VARCHAR(100) NOT NULL,
    customer_number SERIAL PRIMARY KEY,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,
    UNIQUE(uid)
);