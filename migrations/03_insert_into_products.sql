-- +goose Up
INSERT INTO products (name, price, stock) VALUES
('MacBook Pro', 1999.99, 8),
('iPhone 15', 999.00, 15),
('iPad Air', 599.99, 12),
('Apple Watch', 399.00, 20),
('AirPods Pro', 249.99, 30);

-- +goose Down
DELETE FROM products WHERE
(name = 'MacBook Pro' AND  price = 1999.99 AND stock = 8)
OR (name = 'iPhone 15' AND  price = 999.99 AND stock = 15)
OR (name = 'iPad Air', price = 599.99, stock = 12)
OR (name = 'Apple Watch', price = 399.00, stock = 20)
OR (name = 'AirPods Pro', price = 249.99, stock = 30);