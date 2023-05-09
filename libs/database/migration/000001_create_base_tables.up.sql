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

CREATE SCHEMA booking_service;

CREATE TABLE booking_service.bookings (
    booking_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    customer_id TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    shipping_method VARCHAR(50),
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,
    PRIMARY KEY (booking_id)
);

CREATE SCHEMA product_service;

CREATE TABLE product_service.products (
    product_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    rrp INTEGER NOT NULL,
    fit_notes TEXT NOT NULL,
    designer TEXT NOT NULL,
    category VARCHAR(100) NOT NULL,
    shipping_price INTEGER NOT NULL,
    size INTEGER NOT NULL,
    images TEXT[],
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,
    customer_id TEXT NOT NULL,
    PRIMARY KEY (product_id)
);

CREATE SCHEMA payment_service;

CREATE TABLE payment_service.payments (
    payment_id TEXT NOT NULL,
    checkout_session_id TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    booking_id TEXT NOT NULL,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,
    PRIMARY KEY (payment_id)
);