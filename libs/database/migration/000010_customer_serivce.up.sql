CREATE SCHEMA customer_service;

CREATE TABLE customer_service.customers (
    customer_id TEXT NOT NULL,
    email VARCHAR(100) NOT NULL,
    customer_number SERIAL PRIMARY KEY,
    stripe_id VARCHAR(50) UNIQUE,
    stripe_connect_id VARCHAR(50) UNIQUE,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,
    UNIQUE(customer_id)
);