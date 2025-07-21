-- +goose Up
CREATE TABLE IF NOT EXISTS products (
    product_id BIGSERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	price DECIMAL(10,2) NOT NULL,
	stock INTEGER CHECK (stock >= 0)
);

-- +goose Down
DROP TABLE IF EXISTS products;