CREATE TABLE IF NOT EXISTS coffee_shop.orders (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    total DECIMAL(10, 2) NOT NULL,
    order_status ENUM('pending', 'completed', 'cancelled') DEFAULT 'pending',
    created_by VARCHAR(36),
    Updated_by VARCHAR(36),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
    is_deleted BOOLEAN DEFAULT false,
)

CREATE TABLE IF NOT EXISTS order_items (
    id VARCHAR(36) PRIMARY KEY,
    order_id VARCHAR(36) NOT NULL REFERENCES orders(id),
    menu_id VARCHAR(36) NOT NULL REFERENCES menu(id),
    quantity INT NOT NULL,
    price DECIMAL(10, 2) NOT NULL, -- Price of each menu item
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW() ON UPDATE NOW(),
);
