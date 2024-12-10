CREATE TABLE IF NOT EXISTS coffee_shop.users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    email_confirmation_code VARCHAR(255) NULL,
    is_email_confirmed BOOLEAN DEFAULT FALSE,
    email_confirmation_expiry TIMESTAMP NULL,
    role ENUM('customer', 'staff', 'admin') NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    created_at TIMESTAMP DEFAULT current_timestamp(),
    updated_at TIMESTAMP DEFAULT current_timestamp(),
    deleted_at TIMESTAMP
);