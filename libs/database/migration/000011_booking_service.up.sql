CREATE SCHEMA booking_service;

CREATE TABLE booking_service.bookings (
    booking_id TEXT NOT NULL,
    product_id TEXT NOT NULL,
    customer_id TEXT NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL,
    PRIMARY KEY (booking_id),
    CONSTRAINT fk_product FOREIGN KEY (product_id) REFERENCES products(product_id),
    CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES users(uid)
);
