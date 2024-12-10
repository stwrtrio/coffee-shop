CREATE TABLE IF NOT EXISTS coffee_shop.menu (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    category_id VARCHAR(36) REFERENCES categories(id),
    availability BOOLEAN DEFAULT true,
    image_url VARCHAR(255),
    ingredients TEXT,
    preparation_time INT,
    calories INT,
    created_by VARCHAR(36) REFERENCES users(id),
    updated_by VARCHAR(36) REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
    is_deleted BOOLEAN DEFAULT false
);