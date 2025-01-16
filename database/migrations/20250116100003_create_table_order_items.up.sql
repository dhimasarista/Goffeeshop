CREATE TABLE order_items (
    id CHAR(36) NOT NULL DEFAULT (UUID()),
    order_id CHAR(36),
    product_id CHAR(36),
    quantity INT(11) NOT NULL,
    amount INT(11) NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP(),
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP(),
    PRIMARY KEY (id),
    KEY order_id_idx (order_id),
    KEY product_id_idx (product_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);
