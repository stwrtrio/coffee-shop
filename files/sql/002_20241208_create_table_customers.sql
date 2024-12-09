CREATE TABLE IF NOT EXISTS coffee_shop.customers (
    id VARCHAR(36) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp(),
    updated_at TIMESTAMP DEFAULT current_timestamp()
);
