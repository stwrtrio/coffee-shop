CREATE TABLE IF NOT EXISTS coffee_shop.staff_shifts (
    id VARCHAR(36) PRIMARY KEY,
    staff_id VARCHAR(36) REFERENCES staff(id) ON DELETE CASCADE,
    shift_start TIMESTAMP NOT NULL,
    shift_end TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT current_timestamp(),
    updated_at TIMESTAMP DEFAULT current_timestamp()
);
