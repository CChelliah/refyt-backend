CREATE TABLE checkout_sessions (
    checkout_session_id TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    booking_id TEXT NOT NULL,
    created_at TIMESTAMP default now() NOT NULL,
    updated_at TIMESTAMP default now() NOT NULL
);